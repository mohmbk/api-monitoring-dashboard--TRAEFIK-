import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import './login.css'
import { apiUrl } from '../../../config';




function Login() {


    const [email , setemail] = useState('') ;
    const [password , setpassword] = useState('') ;


    const login = async (e : React.MouseEvent) => {
        e.preventDefault();
        try {
            const response = await fetch(`${apiUrl}/login` , {
                method : "POST",
                headers : {
                    "Content-Type" : "application/json"
                },

                body : JSON.stringify({
                    email : email ,
                    password : password ,
                })
            })

            if(!response.ok){
                console.log(await response.text());
                return ;
            }
            alert("login succesfully");

            const data = await response.json();
            localStorage.setItem("token" , data.token);
            window.location.href = "/dashboard" ;

        } catch (error) {
            console.log(error);
        }
    }

  return (
    <>
      <section className='login'>
        <form className='form'>
            <h1>Login</h1>
            <input type="email" placeholder='email' className='inputlog' onChange={(e) => setemail(e.target.value)}/>
            <input type="text" placeholder='password' className='inputlog' onChange={(e) => setpassword(e.target.value)}/>
            <input type="button" value='login !!' className='loginbtn' onClick={login}/>
            <h3>
              Don't have an account? <Link to="/signup">Sign up</Link>
            </h3>
        </form>
      </section>
    </>
  )
}

export default Login
