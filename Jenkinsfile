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
        appEnv = "$ENV"
        appDir = "$appEnv/$BUILD_NUMBER/app"
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
                        error("Error copying repository to $appDir: $err")
                        return
                    }
                }
            }
        }
    }
}