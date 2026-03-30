import http from 'k6/http';
import { randomString } from 'k6/crypto'
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '30s', target: 500 },
        { duration: '1m', target: 500 },
        { duration: '30s', target: 0 },
    ],
};


export default function () {
    const randomUrl = `https://example.com/${Math.random().toString(36).substring(7)}`
    // Создать ссылку
    const res = http.post('http://localhost:8080/shorten',
        JSON.stringify({ url: randomUrl }),
        { headers: { 'Content-Type': 'application/json' } }
    )
    // const id = JSON.parse(res.body).id
    const id = JSON.parse(res.body).id

    for (let i = 0; i < 5; i++) {
        http.get(`http://localhost:8080/${id}`, { redirects: 0 })
        sleep(0.1)
    }
    
    sleep(1)
}