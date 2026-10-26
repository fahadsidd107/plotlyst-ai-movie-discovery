import { Outlet } from "react-router-dom";

const Layout = () => {
    return (
        <main id="main-content">
            <Outlet/>
        </main>
    )
}
export default Layout
