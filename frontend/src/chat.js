import { useState, useEffect } from 'react'

export default function Chat({ nickname }) {
	const [message, setMessage] = useState('')
	const [messages, setMessages] = useState([])
	const [ws, setWs] = useState(null)

	useEffect(() => {
		if (!nickname) return

		const socket = new WebSocket('ws://localhost:8083/ws')

		socket.onopen = () => {
			console.log('WebSocket открыт')
			const welcomeMsg = { name: nickname, content: 'подключился к чату' }
			socket.send(JSON.stringify(welcomeMsg))
			setMessages(prev => [
				...prev,
				`${welcomeMsg.name}: ${welcomeMsg.content}`,
			])
		}

		socket.onmessage = event => {
			console.log('Получено сообщение:', event.data)
			const msg = JSON.parse(event.data)
			setMessages(prev => [...prev, `${msg.name}: ${msg.content}`])
		}

		socket.onerror = error => {
			console.error('Ошибка WebSocket:', error)
		}

		setWs(socket)

		return () => socket.close()
	}, [nickname])

	const sendMessage = () => {
		if (ws && message.trim()) {
			ws.send(JSON.stringify({ name: nickname, content: message }))
			setMessage('')
		}
	}

	return (
		<div className='h-screen flex items-center justify-center bg-gray-100'>
			<div className='bg-white p-6 rounded-2xl shadow-md w-96'>
				<h2 className='text-xl font-semibold mb-4 text-center'>Чат</h2>
				<div className='border p-3 h-64 overflow-y-auto rounded-md bg-gray-50'>
					{messages.map((msg, index) => (
						<div
							key={index}
							className={`p-2 my-1 rounded-md ${
								msg.startsWith(nickname)
									? 'bg-blue-100 text-blue-700 text-right'
									: 'bg-gray-200 text-gray-700'
							}`}
						>
							{msg}
						</div>
					))}
				</div>
				<div className='flex mt-3'>
					<input
						className='border p-2 flex-grow rounded-md'
						placeholder='Введите сообщение...'
						value={message}
						onChange={e => setMessage(e.target.value)}
					/>
					<button
						className='bg-green-500 text-white p-2 ml-2 rounded-md hover:bg-green-600'
						onClick={sendMessage}
					>
						➤
					</button>
				</div>
			</div>
		</div>
	)
}
