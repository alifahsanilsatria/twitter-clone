pipeline{
    agent{
        docker {
            image 'docker:26.0.1-dind-alpine3.19'
        }
    }
    stages{
        stage('build'){
            steps {
                sh "echo 'Hello World (build)'"
            }
        }
    }
    post{
        always{
            echo "echo 'Hello World (always)'"
        }
    }
}