import { Link } from 'react-router-dom'; import LoadingScreen from './LoadingScreen';
const LoadingLab=()=> <div className="loading-lab"><LoadingScreen/><div className="loading-lab__controls"><span>/test · loading laboratory</span><Link to="/">Exit lab</Link><button onClick={()=>window.location.reload()}>Replay</button></div></div>; export default LoadingLab;
