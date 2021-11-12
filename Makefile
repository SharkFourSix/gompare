OUTPUT_NAME = gompare

build:
	go build -o ${OUTPUT_NAME}

install:
	go install

clean:
	go clean

# Tidy up dependencies
tidy:
	go mod tidy

good-test:
	./${OUTPUT_NAME} --template-file=template.csv --input-file=test-good.csv --show-unmatched-cols

bad-test:
	./${OUTPUT_NAME} --template-file=template.csv --input-file=test-bad.csv --show-unmatched-cols