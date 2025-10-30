const baseUrl = import.meta.env.VITE_API_BASE === 'window' ?
	window.location.origin : import.meta.env.VITE_API_BASE;

export const apiUrl = `${baseUrl}/api`;
export const socketUrl = baseUrl.replace(/^http/, 'ws') + '/ws';

const config = Object.freeze({
	"apiUrl": apiUrl,
	"socketUrl": socketUrl
})

export default config;