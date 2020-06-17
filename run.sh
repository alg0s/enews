function check_file_exist {
	echo $1
	if [ -e file ]; 
	then echo "File exists: $1"
	else echo "File does not exist: $1"
	fi
}

function vnnlp {
	MAX_HEAP_SIZE="-Xmx2g"
	SERVER_JAR_FILE="lib/vncorenlp/VnCoreNLPServer.jar"
	NLP_JAR_FILE="lib/vncorenlp/VnCoreNLP-1.1.1.jar"
	ENDPOINT="127.0.0.1"
	ANNOTATORS="wseg,pos,ner,parse"
	PORT="9000"

	check_file_exist $SERVER_JAR_FILE
	check_file_exist $NLP_JAR_FILE

	# cmd line
	java $MAX_HEAP_SIZE -jar $SERVER_JAR_FILE $NLP_JAR_FILE -i $ENDPOINT -p $PORT -a $ANNOTATORS
}

vnnlp
