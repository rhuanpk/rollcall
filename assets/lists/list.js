// let names = ''; document.querySelectorAll('tbody > tr > :nth-child(2)').forEach(e => {const name = e.innerText; name && (names += name+'\n')}); console.log(names)
let names = ''
document.
	querySelectorAll('tbody > tr > :nth-child(2)').
	forEach(e => {
		const name = e.innerText
		name && (names += name + '\n')
	})
console.log(names)
