pipeline {
    agent any

    environment {
        // 基础信息
        APP_NAME      = 'my-go-app'
        DOCKER_IMAGE  = 'sunstarfish/my-go-app'
        DOCKER_REGISTRY = 'https://index.docker.io/v1/'

        // Go 环境
        GOPROXY = 'https://goproxy.cn,direct'
        GOSUMDB = 'off'
        CGO_ENABLED = '0'
    }

    options {
        timestamps()
        disableConcurrentBuilds()
    }

    stages {

        stage('Checkout') {
            steps {
                echo "📥 Checkout 代码"
                checkout scm
            }
        }

        stage('Go Build & Test') {
            agent {
                docker {
                    image 'golang:1.23-alpine'
                    args '-u root:root'
                }
            }
            steps {
                sh '''
                set -e

                echo "🔧 安装基础工具"
                apk add --no-cache git curl

                echo "🐹 Go 环境"
                go version
                go env | grep -E "(GOPROXY|GOSUMDB)"

                echo "📦 下载依赖"
                go mod download
                go mod tidy

                echo "🧪 单元测试"
                go test ./...

                echo "🏗 编译"
                go build -o app
                '''
            }
        }

        stage('Prepare Dockerfile') {
            when {
                expression {
                    return !fileExists('Dockerfile')
                }
            }
            steps {
                sh '''
                echo "📝 生成 Dockerfile"
                cat > Dockerfile << 'EOF'
FROM golang:1.23-alpine AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app

FROM alpine:3.18
RUN apk add --no-cache ca-certificates \
    && addgroup -S app && adduser -S app -G app
WORKDIR /app
COPY --from=builder /app/app .
USER app
CMD ["./app"]
EOF
                '''
            }
        }

        stage('Docker Build') {
            steps {
                sh '''
                echo "🐳 构建 Docker 镜像"
                docker build -t ${DOCKER_IMAGE}:${BUILD_NUMBER} .
                docker tag ${DOCKER_IMAGE}:${BUILD_NUMBER} ${DOCKER_IMAGE}:latest
                '''
            }
        }

        stage('Push Docker Image') {
            steps {
                withDockerRegistry(
                    credentialsId: 'DockerHub',
                    url: "${DOCKER_REGISTRY}"
                ) {
                    sh '''
                    echo "🚀 推送镜像"
                    docker push ${DOCKER_IMAGE}:${BUILD_NUMBER}
                    docker push ${DOCKER_IMAGE}:latest
                    '''
                }
            }
        }

        stage('Deploy') {
            when {
                expression {
                    return fileExists('docker-compose.yml')
                }
            }
            steps {
                sh '''
                echo "🚢 部署应用"
                docker --help
                docker compose down || true
                docker compose pull
                docker compose up -d
                docker ps
                '''
            }
        }
    }

    post {
        success {
            echo "✅ 构建成功：${DOCKER_IMAGE}:${BUILD_NUMBER}"
        }
        failure {
            echo "❌ 构建失败"
        }
        always {
            cleanWs()
        }
    }
}
