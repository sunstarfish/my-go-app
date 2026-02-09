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

        // 宿主机信息
        CONTAINER_PORT = "8000"
        DEPLOY_DIR = "/home/sucre/repos/my-go-app"
        DEPLOY_USER = "fnkf"
        DEPLOY_HOST = "192.168.5.103"
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
                // sshagent(['fnkf']) {
                //     sh """
                //         ssh -o StrictHostKeyChecking=no ${DEPLOY_USER}@${DEPLOY_HOST} '
                //             cd ${DEPLOY_DIR} &&
                //             docker compose up -d ${IMAGE_NAME}
                //         '
                //     """
                // }
                sh '''
                echo "🚢------jenkins安装docker compose---------"
                echo "创建插件目录"
                mkdir -p ~/.docker/cli-plugins
                echo "下载最新版（以 Linux x86_64 为例）"
                curl -SL https://github.com/docker/compose/releases/download/v2.23.0/docker-compose-linux-x86_64 -o ~/.docker/cli-plugins/docker-compose
                echo "添加执行权限"
                chmod +x ~/.docker/cli-plugins/docker-compose

                echo "🚢 部署应用"
                echo "🚢------查看docker版本---------"
                docker --version
                echo "🚢------查看docker详情---------"
                docker info
                echo "🚢------查看docker帮助信息---------"
                docker info
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
