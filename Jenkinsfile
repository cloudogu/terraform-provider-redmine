#!groovy
@Library(['github.com/cloudogu/zalenium-build-lib@v2.1.0']) _
import com.cloudogu.ces.zaleniumbuildlib.*

node('docker') {
    timestamps {
        branch = "${env.BRANCH_NAME}"

        stage('Checkout') {
            checkout scm
        }

        def redmineFilesDir="${WORKSPACE}/redmine-files"
        sh "mkdir -p ${redmineFilesDir}"
        def redmineImage = docker.image('redmine:4.1.1-alpine')

        withDockerNetwork { buildnetwork ->
            redmineImage.withRun("--network ${buildnetwork} " +
                    "-e REDMINE_SECRET_KEY_BASE=supersecretkey " +
                    "-v ${WORKSPACE}/docker-compose/settings.yml:/usr/src/redmine/config/settings.yml " +
                    "-v ${redmineFilesDir}:/usr/src/redmine/files") { redmineContainer ->

                docker.image('golang:1.14.13').inside("--network ${buildnetwork} -e HOME=/tmp") {
                    stage('Build') {
                        make 'clean package checksum'
                        archiveArtifacts 'target/*'
                    }

                    stage('Unit Test') {
                        make 'unit-test'
                        junit allowEmptyResults: true, testResults: 'target/unit-tests/*-tests.xml'
                    }

                    stage('Static Analysis') {
                        make 'static-analysis'
                    }

                    stage('Acceptance Tests') {
                        withEnv(["REDMINE_URL=http://${redmineIP}:3000/",
                                 "REDMINE_CONTAINERNAME=${redmineContainer.id}"
                        ]) {
                            def redmineIP=findIp(redmineContainer)
                            make("wait-for-redmine load-redmine-defaults mark-admin-password-as-changed")
                            make("acceptance-test")
                            archiveArtifacts 'target/acceptance-tests/*.out'
                        }
                    }

                }
            }
        }
        stage('SonarQube') {
            def scannerHome = tool name: 'sonar-scanner', type: 'hudson.plugins.sonar.SonarRunnerInstallation'
            withSonarQubeEnv {
                if (branch == "main") {
                    echo "This branch has been detected as the main branch."
                    sh "${scannerHome}/bin/sonar-scanner"
                } else if (branch == "develop") {
                    echo "This branch has been detected as the develop branch."
                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME} -Dsonar.branch.target=master"
                } else if (env.CHANGE_TARGET) {
                    echo "This branch has been detected as a pull request."
                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.CHANGE_BRANCH}-PR${env.CHANGE_ID} -Dsonar.branch.target=${env.CHANGE_TARGET}"
                } else if (branch.startsWith("feature/")) {
                    echo "This branch has been detected as a feature branch."
                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME} -Dsonar.branch.target=develop"
                }
            }
            timeout(time: 2, unit: 'MINUTES') { // Needed when there is no webhook for example
                def qGate = waitForQualityGate()
                if (qGate.status != 'OK') {
                    unstable("Pipeline unstable due to SonarQube quality gate failure")
                }
            }
        }
    }
}

void make(String goal) {
    sh "make ${goal}"
}

String findIp(container) {
    def containerIP = sh (returnStdout: true, script: "docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${container.id}")
    return containerIP.trim()
}