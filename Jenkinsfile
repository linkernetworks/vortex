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
                withEnv(["GOPATH+AA=${env.WORKSPACE}"]) {
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
    }
}