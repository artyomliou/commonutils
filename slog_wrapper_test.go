package commonutils_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/artyomliou/commonutils"
	"github.com/stretchr/testify/assert"
)

func TestConvertLevelStringToSlogLevel(t *testing.T) {
	type testcase struct {
		input          string
		expectedError  error
		expectedOutput slog.Level
	}
	testcases := []testcase{
		{"debug", nil, slog.LevelDebug},
		{"info", nil, slog.LevelInfo},
		{"warn", nil, slog.LevelWarn},
		{"error", nil, slog.LevelError},

		{"DEBUG", commonutils.ErrInvalidLogLevel, 0},
		{"INFO", commonutils.ErrInvalidLogLevel, 0},
		{"WARN", commonutils.ErrInvalidLogLevel, 0},
		{"ERROR", commonutils.ErrInvalidLogLevel, 0},
	}
	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			slogLevel, err := commonutils.ConvertLevelStringToSlogLevel(tc.input)
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.Equal(t, tc.expectedOutput, slogLevel)
			}
		})
	}
}

func TestUseLoggerFuncs(t *testing.T) {
	type testdata struct {
		logFunc                  func(string, ...any)
		msg                      string
		args                     []any
		expectedKeywordsInOutput []string
	}
	type testcase struct {
		useLoggerFunc func(slog.Level, io.Writer) *slog.Logger
		logLevel      slog.Level
		testData      []testdata
	}
	testcases := []testcase{
		{
			useLoggerFunc: commonutils.UseTextLogger,
			logLevel:      slog.LevelDebug,
			testData: []testdata{
				{
					logFunc: slog.Debug,
					msg:     "test msg",
					args:    []any{slog.String("user", "testuser001")},
					expectedKeywordsInOutput: []string{
						`level=DEBUG`,
						`msg="test msg"`,
						`user=testuser001`,
					},
				},
				{
					logFunc: log.Printf,
					msg:     "test msg",
					args:    []any{"user", "testuser001"},
					expectedKeywordsInOutput: []string{
						`level=DEBUG`,
						`msg="test msg%!(EXTRA string=user, string=testuser001)"`,
					},
				},
			},
		},
		{
			useLoggerFunc: commonutils.UseTextLogger,
			logLevel:      slog.LevelInfo,
			testData: []testdata{
				{
					logFunc: slog.Info,
					msg:     "test msg",
					args:    []any{slog.String("user", "testuser001")},
					expectedKeywordsInOutput: []string{
						`level=INFO`,
						`msg="test msg"`,
						`user=testuser001`,
					},
				},
				{
					logFunc: log.Printf,
					msg:     "test msg",
					args:    []any{"user", "testuser001"},
					expectedKeywordsInOutput: []string{
						`level=INFO`,
						`msg="test msg%!(EXTRA string=user, string=testuser001)"`,
					},
				},
			},
		},
		{
			useLoggerFunc: commonutils.UseTextLogger,
			logLevel:      slog.LevelWarn,
			testData: []testdata{
				{
					logFunc: slog.Warn,
					msg:     "test msg",
					args:    []any{slog.String("user", "testuser001")},
					expectedKeywordsInOutput: []string{
						`level=WARN`,
						`msg="test msg"`,
						`user=testuser001`,
					},
				},
				{
					logFunc: log.Printf,
					msg:     "test msg",
					args:    []any{"user", "testuser001"},
					expectedKeywordsInOutput: []string{
						`level=WARN`,
						`msg="test msg%!(EXTRA string=user, string=testuser001)"`,
					},
				},
			},
		},
		{
			useLoggerFunc: commonutils.UseTextLogger,
			logLevel:      slog.LevelError,
			testData: []testdata{
				{
					logFunc: slog.Error,
					msg:     "test msg",
					args:    []any{slog.String("user", "testuser001")},
					expectedKeywordsInOutput: []string{
						`level=ERROR`,
						`msg="test msg"`,
						`user=testuser001`,
					},
				},
				{
					logFunc: log.Printf,
					msg:     "test msg",
					args:    []any{"user", "testuser001"},
					expectedKeywordsInOutput: []string{
						`level=ERROR`,
						`msg="test msg%!(EXTRA string=user, string=testuser001)"`,
					},
				},
			},
		},
		{
			useLoggerFunc: commonutils.UseJSONLogger,
			logLevel:      slog.LevelDebug,
			testData: []testdata{
				{
					logFunc: slog.Debug,
					msg:     "test msg",
					args:    []any{slog.String("user", "testuser001")},
					expectedKeywordsInOutput: []string{
						`"level":"DEBUG"`,
						`"msg":"test msg"`,
						`"user":"testuser001"`,
					},
				},
				{
					logFunc: log.Printf,
					msg:     "test msg",
					args:    []any{"user", "testuser001"},
					expectedKeywordsInOutput: []string{
						`"level":"DEBUG"`,
						`"msg":"test msg%!(EXTRA string=user, string=testuser001)"`,
					},
				},
			},
		},
		{
			useLoggerFunc: commonutils.UseJSONLogger,
			logLevel:      slog.LevelInfo,
			testData: []testdata{
				{
					logFunc: slog.Info,
					msg:     "test msg",
					args:    []any{slog.String("user", "testuser001")},
					expectedKeywordsInOutput: []string{
						`"level":"INFO"`,
						`"msg":"test msg"`,
						`"user":"testuser001"`,
					},
				},
				{
					logFunc: log.Printf,
					msg:     "test msg",
					args:    []any{"user", "testuser001"},
					expectedKeywordsInOutput: []string{
						`"level":"INFO"`,
						`"msg":"test msg%!(EXTRA string=user, string=testuser001)"`,
					},
				},
			},
		},
		{
			useLoggerFunc: commonutils.UseJSONLogger,
			logLevel:      slog.LevelWarn,
			testData: []testdata{
				{
					logFunc: slog.Warn,
					msg:     "test msg",
					args:    []any{slog.String("user", "testuser001")},
					expectedKeywordsInOutput: []string{
						`"level":"WARN"`,
						`"msg":"test msg"`,
						`"user":"testuser001"`,
					},
				},
				{
					logFunc: log.Printf,
					msg:     "test msg",
					args:    []any{"user", "testuser001"},
					expectedKeywordsInOutput: []string{
						`"level":"WARN"`,
						`"msg":"test msg%!(EXTRA string=user, string=testuser001)"`,
					},
				},
			},
		},
		{
			useLoggerFunc: commonutils.UseJSONLogger,
			logLevel:      slog.LevelError,
			testData: []testdata{
				{
					logFunc: slog.Error,
					msg:     "test msg",
					args:    []any{slog.String("user", "testuser001")},
					expectedKeywordsInOutput: []string{
						`"level":"ERROR"`,
						`"msg":"test msg"`,
						`"user":"testuser001"`,
					},
				},
				{
					logFunc: log.Printf,
					msg:     "test msg",
					args:    []any{"user", "testuser001"},
					expectedKeywordsInOutput: []string{
						`"level":"ERROR"`,
						`"msg":"test msg%!(EXTRA string=user, string=testuser001)"`,
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		for _, data := range tc.testData {
			useLoggerFuncName := runtime.FuncForPC(reflect.ValueOf(tc.useLoggerFunc).Pointer()).Name()
			logLevelName := tc.logLevel.String()
			logFuncName := runtime.FuncForPC(reflect.ValueOf(data.logFunc).Pointer()).Name()
			title := fmt.Sprintf("%s %s %s", useLoggerFuncName, logLevelName, logFuncName)

			t.Run(title, func(t *testing.T) {
				var buf bytes.Buffer
				tc.useLoggerFunc(tc.logLevel, &buf)
				data.logFunc(data.msg, data.args...)

				// examine
				logOutput := buf.String()
				if logOutput == "" {
					t.Fatal("Expected log output but got empty string")
				}
				for _, keyword := range data.expectedKeywordsInOutput {
					if !strings.Contains(logOutput, keyword) {
						t.Errorf("Expected log to contain %s, but got: %s", keyword, logOutput)
					}
				}
			})
		}
	}
}
