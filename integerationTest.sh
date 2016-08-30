#/bin/bash
#Host=http://10.1.235.98:8888
Host=http://54.223.58.0:8888
#Host=http://hub.dataos.io/api
user=panxy3@asiainfo.com
passwd=q


Basic=""
AdminBasic=""
Token=""
AdminToken=""

function getBasic() {

    user=$1
    password=$2
    
    if [ -z "$user" ];then
        echo username null
	exit
    fi
    if [ -z "$password" ];then
        echo password null
	exit
    fi
    
    pw=`echo -n $2 | md5sum`
    tmp=${pw:0:32}
    basic=`echo -n $1:$tmp |base64`
    if [ "${basic}" = "" ];then
        echo "no basic avaliable"
        exit
    fi

    echo $basic
}

function getToken() {

   tmp=$1

   if [ -z "$tmp" ];then
       echo basic null
   exit
   fi

    tokenURL=${Host}/permission/mob
    token=`curl  $tokenURL -H "Authorization: Basic $tmp"  -s`
    token=`echo $token | cut -d \" -f 4`
    
    if [ ${#token} -ne 32 ];then
        echo "no token avaliable"
        exit
    fi

    echo "Token $token"
}

function chkResult() {
    msg=`echo $1 | cut -d "," -f 2 | cut -d ":" -f 2`
    if [ "${msg:1:2}" != "OK" ];then
        echo $1 $2 $3 $4 $5 $6 $7 $8 $9            
    fi
}

Basic=$(getBasic $user $passwd)
echo $Basic
Token=$(getToken $Basic)


Rep=Repository_$RANDOM
Item=Dataitem_$RANDOM
Tag=Tag_$RANDOM
Label=Label_$RANDOM
NewLabel=NewLabel_$RANDOM

echo "1.-----------------------------> 		【拥有者】【新增】Rep        ($Rep)"
result=`curl -X POST ${Host}/repositories/$Rep -d '{"repaccesstype": "public","comment": "中国移动北京终端详情","label": ""}' -H "Authorization: $Token" -s`
chkResult $result

echo "2.-----------------------------> 		【拥有者】【新增】Item       ($Rep/$Item) 	     "
result=`curl -X POST ${Host}/repositories/$Rep/$Item -d '{"repaccesstype": "public","comment": "中国移动北京终端详情","label": {"sys": {"supply_style": "api"},"opt": {},"owner": {},"other": {}},"price":[{"times": 1000,"money": 5,"expire":30},{"times": 10000,"money": 45,"expire":30},{"times":100000,"money": 400,"expire":30}]}' -H "Authorization: $Token" -s`
chkResult $result

echo "3.-----------------------------> 		【拥有者】【新增】Tag        ($Rep/$Item/$Tag) 	     "
result=`curl -X POST ${Host}/repositories/$Rep/$Item/$Tag -d '{"comment":"this is a tag"}' -H "Authorization: $Token" -s`
chkResult $result

##########################################################

echo "4. -----------------------------> 		【用户】查询是否star某个DataItem	     "
result=`curl -X GET ${Host}/star/$Rep/$Item -H "Authorization: $Token" -s`
chkResult $result

echo "5. -----------------------------> 		【用户】更改对一个DataItem的star状态  	     "
result=`curl -X PUT ${Host}/star/$Rep/$Item?star=1 -H "Authorization: $Token" -s`
chkResult $result

echo "6. -----------------------------> 		【用户】更改对一个DataItem的star状态  	     "
result=`curl -X PUT ${Host}/star/$Rep/$Item?star=0 -H "Authorization: $Token" -s`
chkResult $result

echo "7. -----------------------------> 		【任意】返回该DataItem的star量	     "
result=`curl -X GET ${Host}/star_stat/$Rep/$Item -s`
chkResult $result

echo "8. -----------------------------> 		【任意】返回该Drepository的star量	     "
result=`curl -X GET ${Host}/star_stat/$Rep -s`
chkResult $result

##########################################################

echo "26.-----------------------------> 		【拥有者】【删除】Tag 	     "
result=`curl -X DELETE ${Host}/repositories/$Rep/$Item/$Tag -H "Authorization: $Token" -s`
chkResult $result

echo "27.-----------------------------> 		【拥有者】【删除】Item 	     "
result=`curl -X DELETE ${Host}/repositories/$Rep/$Item -H "Authorization: $Token" -s`
chkResult $result

echo "28.-----------------------------> 		【拥有者】【删除】Rep 	     "
result=`curl -X DELETE ${Host}/repositories/$Rep -H "Authorization: $Token" -s`
chkResult $result
