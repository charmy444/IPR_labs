'use client';

import { useState, useEffect, useRef } from 'react';
import axios from 'axios';

const SCROLL_THRESHOLD_PX = 80;

interface User {
  id: number;
  username: string;
  first_name: string;
  last_name: string;
  message_count?: number;
}

interface Message {
  id: number;
  userId: number;
  content: string;
  created_at: string;
  is_read: boolean;
}

interface SupportResponse {
  id: number;
  message_id: number;
  staff_id: number;
  content: string;
  created_at: string;
}

interface MessageWithUser {
  id: number;
  userId: number;
  content: string;
  created_at: string;
  is_read: boolean;
  user: User;
  responses: SupportResponse[];
}

export default function Dashboard() {
  const [users, setUsers] = useState<User[]>([]);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [messages, setMessages] = useState<MessageWithUser[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const [stats, setStats] = useState({
    total_messages: 0,
    unread_messages: 0,
    total_responses: 0,
    total_users: 0,
  });
  const [userHasScrolledUp, setUserHasScrolledUp] = useState(false);
  const messagesContainerRef = useRef<HTMLDivElement>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    fetchUsers();
    fetchStats();
  }, []);

  useEffect(() => {
    if (selectedUser) {
      setUserHasScrolledUp(false);
      fetchMessages(selectedUser.id);
    }
  }, [selectedUser]);

  useEffect(() => {
    if (!userHasScrolledUp) {
      messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages, userHasScrolledUp]);

  const handleMessagesScroll = () => {
    const el = messagesContainerRef.current;
    if (!el) return;
    const distanceFromBottom = el.scrollHeight - el.scrollTop - el.clientHeight;
    setUserHasScrolledUp(distanceFromBottom > SCROLL_THRESHOLD_PX);
  };

  const fetchUsers = async () => {
    try {
      const response = await axios.get('/api/users');
      setUsers(response.data.users || []);
    } catch (error) {
      console.error('Error fetching users:', error);
    }
  };

  const fetchMessages = async (userId: number) => {
    try {
      const response = await axios.get(`/api/users/${userId}/messages`);
      setMessages(response.data.messages || []);
    } catch (error) {
      console.error('Error fetching messages:', error);
    }
  };

  const fetchStats = async () => {
    try {
      const response = await axios.get('/api/stats');
      setStats(
        response.data && typeof response.data === 'object'
          ? {
              total_messages: response.data.total_messages ?? 0,
              unread_messages: response.data.unread_messages ?? 0,
              total_responses: response.data.total_responses ?? 0,
              total_users: response.data.total_users ?? 0,
            }
          : { total_messages: 0, unread_messages: 0, total_responses: 0, total_users: 0 }
      );
    } catch (error) {
      console.error('Error fetching stats:', error);
    }
  };

  const sendMessage = async () => {
    if (!selectedUser || !newMessage.trim()) return;

    try {
      const lastMsg = messages[messages.length - 1];
      await axios.post('/api/responses', {
        message_id: lastMsg?.id,
        response_text: newMessage.trim(),
      });
      setNewMessage('');
      fetchMessages(selectedUser.id);
      fetchStats();
    } catch (error) {
      console.error('Error sending message:', error);
    }
  };

  const markAsRead = async (messageId: number) => {
    if (!selectedUser) return;
    try {
      await axios.post(`/api/messages/${messageId}/read`);
      fetchMessages(selectedUser.id);
      fetchStats();
    } catch (error) {
      console.error('Error marking as read:', error);
    }
  };

  return (
    <div className="container-fluid p-0">
      <div className="row g-0">
        {/* Sidebar */}
        <div className="col-md-3 col-lg-2 sidebar">
          <div className="p-3">
            <h5 className="mb-4">
              <i className="bi bi-robot me-2"></i>
              Support Bot
            </h5>
            <nav className="nav flex-column">
              <a className="nav-link active" href="#">
                <i className="bi bi-chat-dots"></i>
                Сообщения
              </a>
              <a className="nav-link" href="#">
                <i className="bi bi-bar-chart"></i>
                Статистика
              </a>
            </nav>
          </div>
          <div className="sidebar-users-wrap">
            <h6 className="text-white-50 mb-3">Пользователи</h6>
            <div className="user-list">
              {(Array.isArray(users) ? users : []).map((user) => (
                <div
                  key={user.id}
                  className={`user-item ${selectedUser?.id === user.id ? 'active' : ''}`}
                  onClick={() => setSelectedUser(user)}
                  role="button"
                  tabIndex={0}
                  onKeyDown={(e) => e.key === 'Enter' && setSelectedUser(user)}
                >
                  <div className="user-name">
                    {user.first_name} {user.last_name}
                  </div>
                  <div className="user-row">
                    <span className="last-message">@{user.username}</span>
                    {typeof user.message_count === 'number' && (
                      <span className="user-message-count">{user.message_count}</span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Main Content */}
        <div className="col-md-9 col-lg-10 main-content">
          <div className="main-content-inner container-fluid p-4">
            <div className="row g-3 mb-4 stats-row">
              <div className="col-md-3">
                <div className="stats-card">
                  <div className="stats-icon text-primary">
                    <i className="bi bi-chat-dots"></i>
                  </div>
                  <div className="stats-value">{stats.total_messages}</div>
                  <div className="stats-label">Всего сообщений</div>
                </div>
              </div>
              <div className="col-md-3">
                <div className="stats-card">
                  <div className="stats-icon text-warning">
                    <i className="bi bi-envelope"></i>
                  </div>
                  <div className="stats-value">{stats.unread_messages}</div>
                  <div className="stats-label">Непрочитанных</div>
                </div>
              </div>
              <div className="col-md-3">
                <div className="stats-card">
                  <div className="stats-icon text-success">
                    <i className="bi bi-reply"></i>
                  </div>
                  <div className="stats-value">{stats.total_responses}</div>
                  <div className="stats-label">Ответов</div>
                </div>
              </div>
              <div className="col-md-3">
                <div className="stats-card">
                  <div className="stats-icon text-info">
                    <i className="bi bi-people"></i>
                  </div>
                  <div className="stats-value">{stats.total_users}</div>
                  <div className="stats-label">Пользователей</div>
                </div>
              </div>
            </div>

            {selectedUser ? (
              <div className="chat-area-wrap">
                <div className="chat-container">
                <div className="p-3 bg-white border-bottom chat-header">
                  <h5 className="mb-0">
                    <i className="bi bi-person-circle me-2"></i>
                    {selectedUser.first_name} {selectedUser.last_name}
                  </h5>
                  <small className="text-muted">@{selectedUser.username}</small>
                </div>

                <div
                  className="messages-area"
                  ref={messagesContainerRef}
                  onScroll={handleMessagesScroll}
                >
                  {(Array.isArray(messages) ? messages : []).map((msg) => (
                    <div key={msg.id}>
                      <div className="message-bubble user">
                        <div>{msg.content}</div>
                        <div className="message-time">
                          {new Date(msg.created_at).toLocaleTimeString('ru-RU', {
                            hour: '2-digit',
                            minute: '2-digit',
                          })}
                        </div>
                      </div>
                      {(msg.responses || []).map((response) => (
                        <div key={response.id} className="message-bubble staff">
                          <div>{response.content}</div>
                          <div className="message-time">
                            {new Date(response.created_at).toLocaleTimeString('ru-RU', {
                              hour: '2-digit',
                              minute: '2-digit',
                            })}
                          </div>
                        </div>
                      ))}
                    </div>
                  ))}
                  <div ref={messagesEndRef} className="chat-scroll-anchor" aria-hidden="true" />
                </div>

                <div className="input-area">
                  <div className="input-group">
                    <input
                      type="text"
                      className="form-control"
                      placeholder="Введите сообщение..."
                      value={newMessage}
                      onChange={(e) => setNewMessage(e.target.value)}
                      onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
                    />
                    <button className="btn btn-primary" onClick={sendMessage} type="button">
                      <i className="bi bi-send"></i>
                    </button>
                  </div>
                </div>
              </div>
              </div>
            ) : (
              <div className="d-flex align-items-center justify-content-center chat-area-wrap">
                <div className="text-center text-muted">
                  <i className="bi bi-chat-square-text" style={{ fontSize: '4rem' }}></i>
                  <h4 className="mt-3">Выберите пользователя</h4>
                  <p>для начала переписки</p>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
