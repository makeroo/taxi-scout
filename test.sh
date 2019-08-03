COMMAND="$(which $0)"
COMMAND_HOME="$(dirname ${COMMAND})"
PROJECT_HOME=github.com/makeroo/taxi_scout

function run_test {
    TEST_SOURCE="$1"
    TEST_PACKAGE="$(dirname ${TEST_SOURCE})"

    go test ${PROJECT_HOME}/${TEST_PACKAGE}
}

cd "${COMMAND_HOME}"

for test_file in $(find * -name "*_test.go")
do
    run_test "${test_file}"
done
