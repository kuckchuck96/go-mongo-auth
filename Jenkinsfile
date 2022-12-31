pipeline {
    agent any
    tools {
        go 'go-1.19'
    }
    options {
        skipStagesAfterUnstable()
        timeout(time: 1, unit: 'HOURS') 
    }
    parameters {
        choice(name: 'ENV', choices: ['Dev', 'Prod'], description: 'Pick application environment')
    }
    environment {
        appEnv = "$ENV".toLowerCase()
        appDir = "$appEnv/$BUILD_NUMBER/app"
        dockerHost = 'https://docker.io'
        dockerAuth = credentials('docker-hub-auth')
    }
    stages {
        stage('Checkout') {
            steps {
                script {
                    try {
                        sh "mkdir -p $appDir"
                        dir(appDir) {
                            checkout scm
                        }
                    } catch(Exception ex) {
                        error("Error copying repository to $appDir: $ex")
                        return
                    }
                }
            }
        }
        stage('Build') {
            steps {
                script {
                    try {
                        dir(appDir) {
                            sh 'make build'
                        }
                    } catch(Exception ex) {
                        error("Error building project: $ex")
                        return
                    }
                }
            }
        }
        stage('Create & Push Docker Image') {
            steps {
                script {
                    try {
                        dir(appDir) {
                            docker.withRegistry(dockerHost, dockerAuth) {
                                def dockerImage = docker.build("go-mongo-auth:$GIT_COMMIT")
                            }
                        }
                    } catch(Exception ex) {
                        error("Error building or pushing docker image: $ex")
                        return
                    }
                }
            }
        }
    }
}