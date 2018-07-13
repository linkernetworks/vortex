pipeline {
    agent {
        dockerfile {
            dir "src/github.com/linkernetworks/vortex/jenkins"
            args "--privileged --group-add docker"
        }
    }
    post {
        always {
            dir ("src/github.com/linkernetworks/vortex") {
                sh "make clean"
            }
        }
        failure {
            script {
                def message =   "<https://jenkins.linkernetworks.co/job/vortex/job/vortex/|vortex> » " +
                                "<${env.JOB_URL}|${env.BRANCH_NAME}> » " +
                                "<${env.BUILD_URL}|#${env.BUILD_NUMBER}> failed."

                slackSend channel: '#09_jenkins', color: 'danger', message: message

                switch (env.BRANCH_NAME) {
                    case ~/(.+-)?rc(-.+)?/:
                    case ~/develop/:
                    case ~/master/:
                        message += " <!here>"
                        slackSend channel: '#01_vortex', color: 'danger', message: message
                        break
                }
            }
        }
        fixed {
            slackSend channel: '#09_jenkins', color: 'good',
                message:    "<https://jenkins.linkernetworks.co/job/vortex/job/vortex/|vortex> » " +
                            "<${env.JOB_URL}|${env.BRANCH_NAME}> » " +
                            "<${env.BUILD_URL}|#${env.BUILD_NUMBER}> is fixed."
        }
    }
    options {
        timestamps()
        timeout(time: 1, unit: 'HOURS')
        checkoutToSubdirectory('src/github.com/linkernetworks/vortex')
    }
    stages {
        stage('Prepare') {
            steps {
                withEnv(["GOPATH+AA=${env.WORKSPACE}"]) {
                    dir ("src/github.com/linkernetworks/vortex") {
                        sh "make pre-build"
                    }
                }
            }
        }
        stage('Build') {
            steps {
                withEnv(["GOPATH+AA=${env.WORKSPACE}"]) {
                    dir ("src/github.com/linkernetworks/vortex") {
                        sh "make build"
                    }
                }
            }
        }
        stage('Test') {
            steps {
                withEnv(["GOPATH+AA=${env.WORKSPACE}", "TEST_PROMETHEUS=1"]) {
                    dir ("src/github.com/linkernetworks/vortex") {
                        sh "make src.test-coverage 2>&1 | tee >(go-junit-report > report.xml)"
                        junit "report.xml"
                        sh 'gocover-cobertura < build/src/coverage.txt > cobertura.xml'
                        cobertura coberturaReportFile: "cobertura.xml", failNoReports: true, failUnstable: true
                        publishHTML (target: [
                            allowMissing: true,
                            alwaysLinkToLastBuild: true,
                            keepAll: true,
                            reportDir: 'build/src',
                            reportFiles: 'coverage.html',
                            reportName: "GO cover report",
                            reportTitles: "GO cover report",
                            includes: "coverage.html"
                        ])
                    }
                }
            }
        }
        stage("Build Image"){
            when {
                branch 'develop'
            }
            steps {
                script {
                    dir ("src/github.com/linkernetworks/vortex") {
                        docker.build("sdnvortex/vortex" , "--file ./dockerfiles/Dockerfile .")
                    }
                }
            }
        }
        stage("Push Image"){
            when {
                branch 'develop'
            }
            steps {
                script {
                    withCredentials([
                        usernamePassword(
                            credentialsId: 'eb1d8dd2-afd2-49d3-bbef-605de4f664d2',
                            usernameVariable: 'DOCKER_USER',
                            passwordVariable: 'DOCKER_PASS'
                        )
                    ]) {
                        sh 'echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin'
                    }
                    docker.image("sdnvortex/vortex").push("develop")
                }
            }
        }
    }
}
