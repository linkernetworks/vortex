pipeline {
    agent {
        dockerfile {
            dir "jenkins"
            args "--privileged --group-add docker"
        }
    }
    post {
        always {
            sh "make clean"
        }
    }
    stages {
        stage('Prepare') {
            steps {
                sh "go get -u github.com/kardianos/govendor"
                sh "make pre-build"
                sh "docker run -itd -p 27017:27017 --name mongo mongo"
            }
        }
        stage('Build') {
            steps {
                sh "make build"
            }
        }
        stage('Test') {
            steps {
                sh "make src.test-coverage"
            }
        }
    }
}