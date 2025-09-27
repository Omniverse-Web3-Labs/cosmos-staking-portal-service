package server

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

// 判断是不是lambda环境，
// 读环境变量或者配置,这是测试
func InLambda() bool {
	if lambdaTaskRoot := os.Getenv("AWS_LAMBDA"); lambdaTaskRoot != "" {
		return true
	}
	return false
}

// lambda的event处理
func HandlerLambda(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func ServeLambda(engine *gin.Engine) {
	ginLambda = ginadapter.New(engine)
	lambda.Start(HandlerLambda)
}
