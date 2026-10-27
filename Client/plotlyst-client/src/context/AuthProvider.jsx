import {createContext, useState,useEffect, useCallback} from 'react';
import { authApi } from '../api/step4';

const AuthContext = createContext({});

export const AuthProvider = ({children}) => {
    const [auth, setAuth] = useState(null);
    const [loading, setLoading] = useState(true);
    const hydrate = useCallback(async () => {
        try { setAuth(await authApi.session()); }
        catch { setAuth(null); }
        finally { setLoading(false); }
    }, []);
    useEffect(() => { hydrate(); }, [hydrate]);

    return (
        <AuthContext.Provider value = {{auth,setAuth,loading,hydrate}}>
            {children}
        </AuthContext.Provider>
    )
}
export default AuthContext;
