#set HOST https://www.baidu.com
#set POSTDATA username=a,password=a
#req POST /login/
#req GET /
#assert STATUS 200
#echo $RESP

set HOST http://203.195.180.22
set HEADER Qbtoken:c8a7d157132d8fe0eb5b78aacadf9f4c3a3ce11b#,Uuid:IMEI_e9f544e8ded394a4d5d2eaba6acb8298

for ${typ} in [suggest,txt]

echo ${typ}
req GET /article/list/${typ}
#assert STATUS 200
getvalue ${err} err
#assert ${err} 0
#echo $RESP
getarray ${arr} id
#echo ${arr}

set ${ids} ${arr}
echo ${ids}
for ${aid} in ${ids}
echo ${aid}
req GET /article/${aid}
#assert STATUS 200
#echo $RESP
getvalue ${err} err
#assert ${err} 0
endfor

endfor
