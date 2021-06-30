# bazel coverage --define version=1.0 --combined_report=lcov --runs_per_test=1 --flaky_test_attempts=3 -- //:paloma-config-service_test
# find bazel-testlogs/  -type f -name coverage.dat -exec cat {} + > tmp.dat
# sort tmp.dat | uniq > coverage_report.dat
# export SONAR_SCANNER_OPTS="-Djavax.net.ssl.trustStore=/Library/Java/JavaVirtualMachines/jdk1.8.0_231.jdk/Contents/Home/jre/lib/security/cacerts -Djavax.net.ssl.trustStorePassword=changeit"
export SONAR_SCANNER_OPTS="-Djavax.net.ssl.trustStore=/usr/lib/jvm/java-11-openjdk-amd64/lib/security/cacerts -Djavax.net.ssl.trustStorePassword=changeit"
# go vet src/*.go
bazel coverage --define version=$BUILD_NUMBER --combined_report=lcov --runs_per_test=1 --flaky_test_attempts=3 -- //:paloma-config-service_test
#~/softwares/sonar-scanner-4.6.1.2450-macosx/bin/sonar-scanner
/opt/sonar-scanner/bin/sonar-scanner -Dsonar.login=${SONAR}  -Dsonar.branch.name=${BRANCH_NAME}