
interface IPortal {
  id: string;
  img: string;
  name: string;
  redirectUrl: string;
  desc: string;
  count: number;
  //updateCount: (c: number) => void;
  //redirect: (url: string) => void;
}

export default IPortal