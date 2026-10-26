import { useCallback, useEffect, useState } from 'react';

export default function useCatalog(loader, dependencies = []) {
  const [state, setState] = useState({ loading: true, data: null, meta: null, error: null });
  const load = useCallback(async () => {
    setState(current => ({ ...current, loading: true, error: null }));
    try {
      const response = await loader();
      setState({ loading: false, data: response.data.data, meta: response.data.meta, error: null });
    } catch (error) {
      setState({ loading: false, data: null, meta: null, error: error.response?.data?.error?.message || 'The catalogue signal dropped. Please try again.' });
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, dependencies);
  useEffect(() => { load(); }, [load]);
  return { ...state, retry: load };
}
