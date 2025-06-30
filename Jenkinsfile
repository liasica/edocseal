pipeline {
    agent any

    environment {
        CURRENT_BRANCH = "${env.GIT_BRANCH.replace('origin/', '')}"
        TAG = "${env.CURRENT_BRANCH == 'master' ? 'prod' : 'dev'}"
        HOST = "${env.CURRENT_BRANCH == 'master' ? '172.16.1.10' : '172.16.1.11'}"
        TOKEN = "${env.CURRENT_BRANCH == 'master' ? 'LP6LAA1ricmZ66CmoQvJsZQvCcggfQG0RqDeBwqTU7esVP2F5XYx1a0Mze46v8K6' : 'W3LmegT1Ky5HkaWXywcrKvQs3oG2ojc7YiornB1v2fYpGHJfpm5vEomsi2HByUd1'}"
    }

    stages {
        stage('Diagnosis') {
            steps {
                //checkout scm

                echo "----------------------------------------"
                echo "| 当前分支变量列表："
                echo "| CURRENT_BRANCH: ${env.CURRENT_BRANCH}"
                echo "| TAG: ${env.TAG}"
                echo "| HOST: ${env.HOST}"
                echo "----------------------------------------"
            }
        }

         stage('Go Build') {
            steps {
                sh """
                    go mod download
                    GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -tags=jsoniter,poll_opt -gcflags "all=-N -l" -o build/release/edocseal cmd/edocseal/main.go
                """
            }
        }

        stage('Docker Image') {
            steps {
                script {
                    docker.withRegistry('https://harbor.liasica.com', 'harbor-bot') {
                        def dockerImage = docker.build("harbor.liasica.com/auroraride/edocseal:${env.TAG}", "-f Dockerfile .")
                        dockerImage.push()
                        dockerImage.push("${env.BUILD_ID}")
                    }
                }
            }
        }

        stage('测试') {
            when {
                expression {
                    env.CURRENT_BRANCH == 'development'
                }
            }
            steps {
                script {
                    server_deploy(env.HOST, env.TAG)
                }
            }
        }

        stage('正式') {
            when {
                expression {
                    env.CURRENT_BRANCH == 'master'
                }
            }
            steps {
                script {
                    // 拉取镜像
                    server_deploy(env.HOST, env.TAG)
                }
            }
        }
    }
}

def server_deploy(String host, String tag) {
    def server = [:]
    server.host = host
    server.name = "AUR-${tag}"
    server.allowAnyHosts = true
    withCredentials([sshUserPrivateKey(credentialsId: 'jenkins-root', keyFileVariable: 'identity', passphraseVariable: '', usernameVariable: 'userName')]) {
        server.user = userName
        server.identityFile = identity
        sshCommand remote: server, command: """
            cd /var/www
            docker pull harbor.liasica.com/auroraride/edocseal:${tag}
            docker compose stop edocseal
            docker compose rm -f edocseal
            docker compose up edocseal -d
            docker image prune -f -a
            docker container prune -f
            docker volume prune -f -a
        """
    }
}
