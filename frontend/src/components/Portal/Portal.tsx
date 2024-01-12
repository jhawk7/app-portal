import IPortal from "../../interfaces/Portal";

interface PortalProps {
  portal: IPortal
  updateCount: (id: string, c: number) => void;
  handleRedirect: (url: string) => void;
}

const Portal = ({portal, updateCount, handleRedirect}: PortalProps) => {
//highlight card on hover
  return (
    <>
      <div className="card" 
        onClick={() => {
          updateCount(portal.id, portal.count); 
          handleRedirect(portal.redirectUrl)
        }}>
        <img src={portal.img} className="card-img-top" alt={portal.name + ' image'}></img>
        <div className="card-body">
          <h5 className="card-title">{portal.name}</h5>
          <p className="card-text">{portal.desc}</p>
          {/* <a href={props.redirectUrl} className="btn btn-primary">Re</a> */}
        </div>
    </div>
    </>
  )
}

export default Portal