# blockchain-on-api
api on blockchain project 


to create new book api

http://localhost:8080\new    post request

put body on and to json file 



{
    "title":"simplebook",

    "author":"rezahayati",

    "isbn":"101010",

    "publishdate":"2020,1,6"
}


response well be

{
    "Id": "f1b4a5f8228da4c8960b37b8b87b5bba",
    
    "Title": "simplebook",
    
    "Author": "rezahayati",
    
    "PublishDate": "2020,1,6",
    
    "isbn": "101010"
}



then you get book id and checkout date that reposnsed to you and paste in blockpost route to post a block
http://localhost:8080 post request 



{
   "bookid": "f1b4a5f8228da4c8960b37b8b87b5bba",
   
   "user":"hayati",
   
   "checkoutdate":"2020,1,6"
}



and last step you can get your blockchain api that created by using 
http://localhost:8080   put on get request


and response well be 


   {
        "Pos": 4,
        
        "Data": {
        
            "BookId": "",
            
            "User": "",
            
            "CheckOutDate": "",
            
            "IsGenesis": false
            
        },
    
        "TimeStamp": "2022-11-16 10:20:38.6469125 +0330 +0330 m=+553.980378801",
        
        "Hash": "22bb0c5e354c60853911d7c169af3cca110adc592797798fc9df788437a016b9",
        
        "PrevHash": "73695c47a53758b5afcabbe175f0cb88f0f5af2611c1cb1d9f0141a15847bba8"
         }
         ]
