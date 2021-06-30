pipeline {
  agent {label 'paloma'}
  environment {
    DOCKER_CREDS = credentials('PitArtUser')   
  }
  parameters {    
    booleanParam(name: 'K8S_OBJ_ANALYSIS', defaultValue: false,  description: 'Enables kubernetes object analysis')    
    booleanParam(name: 'STATIC_ANALYSIS', defaultValue: false, description: 'Enables sonar scan')    
    booleanParam(name: 'UPLOAD_DOCKER_IMAGE', defaultValue: false, description: 'Upload to artifactory')    
  }
  stages {
    //  stage('Initialization'){
    //     steps {
    //         // buildName('#${BUILD_NUMBER} (${ENV,var="BRANCH"}, ${BUILD_CAUSE})')
    //         buildName('#${ENV, var="BRANCH_NAME"}')
    //     } // steps
    // } // stage Initialization
    stage('build & test') {
      parallel {
        stage('build') {
          steps {
            sh '''
            bazel build //:paloma-config-service
            '''
          }
        }

        stage('unit test') {
          steps {
            sh '''bazel test //:paloma-config-service_test
'''
          }
        }

      }
    }

     stage('sonar') {
        when {
          expression { params.get("STATIC_ANALYSIS") == true }
        }
        steps {
            withCredentials([string(credentialsId: 'SONAR_SCANNER_TOKEN', variable:'SONAR')]) {
                sh '''sh sonar-coverage.sh
                '''
            }        
            // timeout(time: 1, unit: 'HOURS') {
            //   waitForQualityGate abortPipeline: true
            // }
        }
    }
    
    
    stage('static analysis of kubernetes objects') {
        when {
          expression { params.get("K8S_OBJ_ANALYSIS") == true }
        }
        steps {        
          sh 'docker run -v $(pwd):/project zegl/kube-score:v1.10.0 score *.yaml'               
        }
    }
    

    stage('create docker image') {
      steps {
        sh '''
        bazel build config-service-image
        bazel run config-service-image
        '''
      }
    }

    stage('create infrastructure') {
      steps {
        sh '''
        #cd terraform
# terraform init
# terraform apply -auto-approve
k3d cluster delete paloma
k3d cluster create paloma'''
      }
    }

    stage('deploy application') {
      steps {
        sh '''k3d image import -c paloma bazel:config-service-image
kubectl apply -f deployments/app-deployment.yml
kubectl wait --for=condition=available --timeout=2m deployment/config-service
sleep 30s
'''
      }
    }

    stage('component testing') {
      steps {
        sh '''
        cd tests
        sh resources/test-component.sh
kubectl wait --for=condition=complete --timeout=5m job/component-test            
kubectl logs -l type=component-test
SUCCESS=$(kubectl get job component-test -o jsonpath=\'{.status.succeeded}\')
if [ $SUCCESS != \'1\' ]; then exit 1; fi
echo "Component test succesful"
'''
      }
    }

    stage('destroy infrastructure') {
      steps {
        // script {
        //   input('Do you want to destroy infrastructure?')
        // }

        sh '''# cd terraform
# terraform destroy -auto-approve
k3d cluster delete paloma
'''
      }
    }

    stage('upload to artifactory') {
       when {
          expression { params.get("UPLOAD_DOCKER_IMAGE") == true }
        }
        steps {
          sh '''
          docker login -u $DOCKER_CREDS_USR -p $DOCKER_CREDS_PSW https://artifactory.dc.pa.sophos:6555/v2/
          bazel run --define version=v1.0.5 --define docker_registry=artifactory.dc.pa.sophos:6555/artifactory/paloma-docker-dev --define repository=jaydatt/config-service config-service-push 
          '''
       }
    }

  }
}