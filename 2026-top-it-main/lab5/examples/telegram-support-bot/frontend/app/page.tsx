'use client';

import Link from 'next/link';

export default function Home() {
  return (
    <div className="d-flex align-items-center justify-content-center min-vh-100">
      <div className="hero-card text-center">
        <i className="bi bi-robot hero-icon"></i>
        <h1 className="mb-4">Telegram Support Bot</h1>
        <p className="lead text-muted mb-4">
          Система технической поддержки на базе Telegram-бота
        </p>
        <div className="d-grid gap-2">
          <Link href="/dashboard" className="btn btn-primary btn-lg">
            <i className="bi bi-speedometer2 me-2"></i>
            Перейти к панели управления
          </Link>
        </div>
      </div>
    </div>
  );
}
