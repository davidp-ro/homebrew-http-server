package server

import (
	"testing"
)

type TestPath struct {
	MatcherPath         string
	MatchingReqPaths    []string
	NonMatchingReqPaths []string
	// ExpectedParams will be tested against the params from MatchingReqPaths[0]
	ExpectedParams map[string]string
}

func TestExtractPathParams(t *testing.T) {
	cases := []TestPath{
		{
			MatcherPath: "/$path1/$path2/test",
			MatchingReqPaths: []string{
				"/hello/world/test",
				"/hello/123/test",
				"/test/test/test",
			},
			NonMatchingReqPaths: []string{
				"/hello/world",
				"/hello/world/test/extra",
			},
			ExpectedParams: map[string]string{
				"$path1": "hello",
				"$path2": "world",
				"test":   "test",
			},
		},
		{
			MatcherPath: "/$path1/test/$path2",
			MatchingReqPaths: []string{
				"/hello/test/world",
				"/test/test/world",
				"/test/test/123",
			},
			NonMatchingReqPaths: []string{
				"/hello/world",
				"/hello/world/test",
				"/hello/test/world/extra",
			},
			ExpectedParams: map[string]string{
				"$path1": "hello",
				"$path2": "world",
				"test":   "test",
			},
		},
	}

	for i, currentCase := range cases {
		// Matching request paths
		for j, matchingReqPath := range currentCase.MatchingReqPaths {
			pathParams := extractPathParams(currentCase.MatcherPath, matchingReqPath)
			if pathParams == nil {
				t.Errorf("Expected path params to be not nil for case %d.%d", i, j)
			}

			if j == 0 {
				// Check expected params for the first matching path
				for key, value := range currentCase.ExpectedParams {
					if pathParams[key] != value {
						t.Errorf("Expected path param %s to be %s, got %s", key, value, pathParams[key])
					}
				}

				if checkPathParams(pathParams) == false {
					t.Errorf("Expected checkPathParams to be true for case %d.%d", i, j)
				}
			}
		}

		// Non-matching request paths
		for j, nonMatchingReqPath := range currentCase.NonMatchingReqPaths {
			pathParams := extractPathParams(currentCase.MatcherPath, nonMatchingReqPath)

			if checkPathParams(pathParams) == true {
				t.Errorf("Expected checkPathParams to be false for case %d.%d", i, j)
			}
		}
	}
}
