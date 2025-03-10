import { useState } from 'react'
import Chat from './chat'

export default function App() {
	const [nickname, setNickname] = useState('')
	const [inputName, setInputName] = useState('')

	return (
		<div className='h-screen flex items-center justify-center bg-gray-100'>
			{!nickname ? (
				<div className='bg-white p-6 rounded-2xl shadow-md w-80 text-center'>
					<h2 className='text-xl font-semibold mb-4'>Введите ник</h2>
					<input
						className='border p-2 w-full rounded-md text-center'
						placeholder='Ваш ник'
						value={inputName}
						onChange={e => setInputName(e.target.value)}
					/>
					<button
						className='bg-blue-500 text-white p-2 mt-3 w-full rounded-md hover:bg-blue-600'
						onClick={() => setNickname(inputName.trim())}
					>
						Войти в чат
					</button>
				</div>
			) : (
				<Chat nickname={nickname} />
			)}
		</div>
	)
}
