package utils

import (
	"fmt"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/gin-gonic/gin"
)

type CtxKey string

func (c CtxKey) String() string {
	return string(c)
}

// func fillStruct[T any](input interface{}) (T, error) {
// 	var result T

// 	// Преобразуем входные данные в JSON и заполним целевую структуру
// 	data, err := json.Marshal(input)
// 	if err != nil {
// 		return result, fmt.Errorf("ошибка при преобразовании в JSON: %v", err)
// 	}

// 	err = json.Unmarshal(data, &result)
// 	if err != nil {
// 		return result, fmt.Errorf("ошибка при заполнении структуры: %v", err)
// 	}

// 	return result, nil
// }

func GetDataFromContext[T any](ctx *gin.Context, key string) (T, error) {
	contextData, exists := ctx.Get(CtxKey(key).String())
	var result T
	if !exists {
		return result, domain.ErrNoUserInContext
	}

	fmt.Printf("contextData %v\n", contextData)
	if contextData == nil {
		return result, domain.ErrNoUserInContext
	}
	data, ok := contextData.(T)
	if !ok {
		return result, domain.ErrNoUserInContext
	}
	return data, nil
}

// func GetUserFromContext(ctx *gin.Context) (domain.User, error) {
// 	contextUser, exists := ctx.Get(CtxKey("user").String())
// 	if !exists {
// 		return domain.User{}, domain.ErrNoUserInContext
// 	}
// 	// contextUser := ctx.Value("user")

// 	fmt.Printf("contextUser %v\n", contextUser)
// 	if contextUser == nil {
// 		return domain.User{}, domain.ErrNoUserInContext
// 	}
// 	user, ok := contextUser.(domain.User)
// 	if !ok {
// 		return domain.User{}, domain.ErrNoUserInContext
// 	}
// 	return user, nil
// }
