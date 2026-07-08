package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestVerifyJWT(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()
	tokenSecret1 := "secretToken1"
	tokenSecret2 := "secretToken2"
	expiredTime1 := time.Duration(100) * time.Second
	expiredTime2 := time.Duration(200) * time.Second
	tokenString1, _ := MakeJWT(userID1, tokenSecret1, expiredTime1)
	tokenString2, _ := MakeJWT(userID2, tokenSecret2, expiredTime2)

	tests := []struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		expiredTime time.Duration
		tokenString string
		wantErr     bool
	}{
		{
			name:        "Valid Token 1",
			userID:      userID1,
			tokenSecret: tokenSecret1,
			expiredTime: expiredTime1,
			tokenString: tokenString1,
			wantErr:     false,
		},
		{
			name:        "Valid Token 2",
			userID:      userID2,
			tokenSecret: tokenSecret2,
			expiredTime: expiredTime2,
			tokenString: tokenString2,
			wantErr:     false,
		},
		{
			name:        "Invalid Secret (Token 1 verified with Secret 2)",
			userID:      uuid.Nil,
			tokenSecret: tokenSecret2,
			expiredTime: expiredTime1,
			tokenString: tokenString1,
			wantErr:     true,
		},
		{
			name:        "Malformed Token String",
			userID:      uuid.Nil,
			tokenSecret: tokenSecret1,
			expiredTime: 0,
			tokenString: "this.is.not.a.valid.jwt.token",
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userID, err := VerifyJWT(tc.tokenString, tc.tokenSecret)
			if tc.wantErr {
				if err == nil {
					t.Error("failed to test, expected: error, result: nil")
				}
				return
			}
			if err != nil {
				t.Errorf("failed to test, expected: success, result: %v", err)
			}
			if userID != tc.userID {
				t.Errorf("failed to test, expected: %v, result: %v", tc.userID, userID)
			}

		})
	}

}
