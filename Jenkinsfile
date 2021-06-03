import groovy.json.JsonSlurper
pipeline {
    agent { docker { image 'python:3.5.1' } }
    stages {
        stage('build') {
            steps {
                sh 'python --version'
            }
        }
        stage('Check Sonar') {
            steps {
                withMaven(maven: 'maven-3.3.9') {
                    withSonarQubeEnv('sonar-6') {
                        sh 'mvn clean install sonar:sonar -Dsonar.scm.disabled=true -Dsonar.host.url=$SONAR_HOST_URL'
                    }
                }
                timeout(1) {
                    waitUntil {
                        script {
                            fileExists('target/sonar/report-task.txt')
                        }
                    }
                    waitUntil {
                        script {
                            def taskId = readFile('target/sonar/report-task.txt').split("\n")[3].split("=")[1]
                            def task_response = httpRequest "https://localhost:9000/api/ce/task?id=${taskId}"
                            def task_data = new JsonSlurperClassic().parseText(task_response.content)
                            return (task_data.task.status.equals("SUCCESS"))
                        }
                    }
                }
                script {
                    def response = httpRequest "https://localhost:9000/api/qualitygates/project_status?projectKey=XXX"
                    def data = new JsonSlurperClassic().parseText(response.content)
                    if (data.projectStatus.status == "ERROR") {
                        error("Sonar Quality Gate not met. Check https://localhost:9000/overview?id=XXX")
                    }
                }
            }
        }
    }
}
