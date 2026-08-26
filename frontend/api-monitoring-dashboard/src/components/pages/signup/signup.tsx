import React, { useState } from 'react'
import './signup.css'
import { apiUrl } from '../../../config';


function Signup() {

    interface signup {
        username : string ;
        email : string ;
        password : string ;
    }


    const [name , setname] = useState('');
    const [ email , setemail] = useState('');
    const [password , setpassword] = useState('');

    const signup = async (e : React.MouseEvent) => {
        e.preventDefault();
        try {
            const response = await fetch(`${apiUrl}/signup` , {
                method : "post" ,
                headers : {
                    "Content-Type" : "application/json"
                },

                body : JSON.stringify({
                    username : name ,
                    email : email ,
                    password : password ,
                })

            })

            if(!response.ok){
                console.log(await response.text());
                return ;
            }

            alert("signup succesfully");
            window.location.href = "/login" ;
        } catch (error) {

        }
    }

  return (
    <>
      <section className='signup'>
        <form className='inpdiv'>
            <h1>signup</h1>
            <input type="text" placeholder='name' className='inputlog' onChange={(e) => setname(e.target.value)}/>
            <input type="email" placeholder='email' className='inputlog' onChange={(e) => setemail(e.target.value)} />
            <input type="text" placeholder='password' className='inputlog' onChange={(e) => setpassword(e.target.value)}/>
            <input type="button" value='sign up !!' className='loginbtn' onClick={signup}/>
        </form>
      </section>
    </>
  )
}

export default Signup
