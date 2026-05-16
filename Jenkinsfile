pipeline {
    agent any

    tools {
        go '1.24.3'
    }

    stages {
        stage('build') {
            steps {
                sh 'go mod tidy'
                sh 'go -C . build -o cli cmd/cli/main.go'
                sh 'go -C . build -o scanner cmd/scanner/main.go'
                sh 'go -C . build -o web cmd/web/main.go'
            }
        }
    }

    post {
        success {
            emailext(
                        to: 'mamcer@protonmail.com',
                        subject: "SUCCESS: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]'",
                        body: """<p>Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' succeeded.</p><p>Check console output at <a href='${env.BUILD_URL}'>${env.BUILD_URL}</a></p>""",
                        mimeType: 'text/html'
                    )
        }

        failure {
            emailext(
                        to: 'mamcer@protonmail.com',
                        subject: "FAILURE: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]'",
                        body: """<p>Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' failed.</p><p>Check console output at <a href='${env.BUILD_URL}'>${env.BUILD_URL}</a></p>""",
                        mimeType: 'text/html'
                    )
        }

        always {
            cleanWs()
        }
    }
}
