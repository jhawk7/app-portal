import { useEffect, useState } from 'react'
import IPortal from './interfaces/Portal';
import Portal from './components/Portal/Portal';

const App = () => {
  const [portals, setPortals] = useState<IPortal[]>([]);

  const getPortals = () => {
    fetch("/portals", {
      method: "GET",
      headers: {
      "Content-Type": "application/json",
      },
    })
    .then((response) => response.json())
    .then((res) => {
      const data: IPortal[] = JSON.parse(res)
      console.log(data)
      setPortals(data);
    })
    .catch(e => {
      console.log(e);
    });
  };

  const handleRedirect = (url: string) => {
    location.href = url;
  }

  const updateCount = (id: string, count: number) => {
    //PATCH api call to update count
    console.log(id, count);
  }

  useEffect(getPortals, [])

  return (
    <>
      <div className="container text-center">
        <div className="row row-cols-3">
          {portals.map((portal: IPortal, index: number) => {
            return <div className="col">
              <Portal key={index} portal={portal} handleRedirect={handleRedirect} updateCount={updateCount}/>
            </div>
          })}
        </div>
      </div>
    </>
  )
}

export default App
