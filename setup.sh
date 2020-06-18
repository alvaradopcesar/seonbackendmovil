#!/bin/bash

echo $ENVIRONMENT

if [ "$ENVIRONMENT" = "ci" ]; then
    echo "configurado para : ci"
    cp ProductService2-ci.yml ProductService2.yml
    cp FirebaseServiceAccount-ci.json FirebaseServiceAccount.json
elif [ "$ENVIRONMENT" = "qa" ]; then
    echo "configurado para : qa"
    cp ProductService2-qa.yml ProductService2.yml
    cp FirebaseServiceAccount-qa.json FirebaseServiceAccount.json    
elif [ "$ENVIRONMENT" = "prd" ]; then
    echo "configurado para : prd"
    cp ProductService2-prd.yml ProductService2.yml
    cp FirebaseServiceAccount-prd.json FirebaseServiceAccount.json        
fi

./ProductService2 >> productservive2.log
