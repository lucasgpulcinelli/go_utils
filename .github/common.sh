
project=$(echo $GITHUB_REF_NAME | sed 's/\/.\+//')
version=$(echo $GITHUB_REF_NAME | sed 's/.\+\///')