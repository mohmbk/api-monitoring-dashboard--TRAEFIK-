import React, { useState , useEffect } from 'react'
import './dashboard.css'
import { useNavigate } from 'react-router-dom';
import { apiUrl } from '../../../config';


function Dashboard() {

  interface Api {
  id: string;
  userId: string;
  name: string;
  url: string;
  lastStatus: string;
  lastStatusCode: number;
  lastResponseTime: number;
  lastCheckedAt: string;
  createdAt: string;
}

interface addapi {
  name : string ;
  url : string ;
}


const [name , setname] = useState('');
const [url , seturl] = useState('');

const [apis , setapis] = useState<Api[]>([]) ;

 useEffect(() => {

        async function fetchApis() {
          const token = localStorage.getItem("token");
          const response = await fetch(`${apiUrl}/dashboard` , {
            method : "GET",
            headers : ({
              "Authorization": `Bearer ${token}`
            })
          })

          if(!response.ok){
            alert(await response.text());
            return;
          }

          const data = await response.json();
          setapis(data ?? []);
        }

        fetchApis();
        const interval = setInterval(() => {
          fetchApis();
        }, 60000);
        return () => clearInterval(interval);

    }, []);

    const addapi = async (e : React.MouseEvent) => {
      e.preventDefault();
      try {
        const token = localStorage.getItem("token");
        const response = await fetch(`${apiUrl}/dashboard` , {
          method : "POST" ,
          headers : {
            "Content-Type" : "application/json" ,
            "Authorization": `Bearer ${token}`
          },

          body : JSON.stringify({
            name : name ,
            url : url ,
          })

        })

        if(!response.ok){
          alert(await response.text());
          return ;
        }

        alert("api created");
        setname("");
        seturl("");
      } catch (error) {
        console.log(error);

      }
    }


    async function deleteapi(id : string) {

      try {
        const token = localStorage.getItem("token");
        const resp = await fetch(`${apiUrl}/dashboard/${id}` , {
          method : "DELETE" ,
          headers : {
            "Authorization": `Bearer ${token}`
          },

        });

        if (!resp.ok){
          alert(await resp.text());
          return ;
        }
      } catch (error) {
        console.log(error);
      }
    }

    const navigate = useNavigate();

    function viewHistory(id: string) {
      navigate(`/dashboard/api/${id}/history`);
    }



  return (
    <>
    <br /><br /><br /><br />
      <section className='apis'>
        {apis.map((api) => (
          <div className="api" key={api.id}>
            <div className='left'>
              <h3>name : {api.name}</h3>

              <p>
                <strong>URL :</strong> {api.url}
              </p>

              <p>
                <strong>Status :</strong> {api.lastStatus}
              </p>

              <p>
                <strong>Code :</strong> {api.lastStatusCode}
              </p>

              <p>
                <strong>Response Time :</strong> {api.lastResponseTime} ms
              </p>

              <p>
                <strong>Last Checked :</strong> {api.lastCheckedAt}
              </p>
            </div>

            <div className='right'>
              <button className='mbtn' onClick={() => deleteapi(api.id)}>delete api</button>
              <button className='mbtn' onClick={() => viewHistory(api.id) }>view history</button>
            </div>


          </div>


        ))}


      </section>
        <br /><br /><br /><br />
      <section className='container'>
        <div className='addapi'>
          <h1>add api !!</h1>
          <input type="text" placeholder='enter name' className='inpapi' onChange={(e) => setname(e.target.value)} />
          <input type="text" placeholder='enter url' className='inpapi' onChange={(e) => seturl(e.target.value)}/>
          <input type="button" value='add api'className='apibtn' onClick={addapi}/>
        </div>
      </section>
    </>
  )
}

export default Dashboard
