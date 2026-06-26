import "bulma/css/bulma.min.css";
import "../../styles/viewProfile.css";
import { useCon } from "../../hooks/useCon";
import { Link, useNavigate } from "react-router-dom";
import { useState } from "react";
import UserHeader from "../UserHeader";
import "../../styles/UserHeader.css";

const API = import.meta.env.VITE_API_BASE_URL;

export default function ConnectionsForm() {
  const navigate = useNavigate();
  const { loading, error, data } = useCon();
  const [chatLoading, setChatLoading] = useState<string | null>(null);
  const [disconnecting, setDisconnecting] = useState<string | null>(null);

  const openChat = async (userId: string) => {
    setChatLoading(userId);
    try {
      const token = localStorage.getItem('token');
      const res = await fetch(`${API}/api/chats`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      const chatsData = await res.json();
      const chat = chatsData.chats?.find((c: any) => 
        c.user1_id === userId || c.user2_id === userId
      );
      if (chat) {
        navigate('/chats', { state: { selectedChatId: chat.id } });
      } else {
        alert('Chat not found. Please try again.');
      }
    } catch (err) {
      console.log(err)    
      alert('Failed to open chat');
    } finally {
      setChatLoading(null);
    }
  };

const handleDisconnect = async (userId: string) => {
    
  if (!confirm('Are you sure you want to disconnect? This will remove the connection and chat history.')) {
    console.log('User cancelled');
    return;
  }

  setDisconnecting(userId);
  
  try {
    const token = localStorage.getItem('token');
    const url = `${API}/api/disconnect`;
        
    const res = await fetch(url, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ target_user_id: userId })
    });
    
   
    
    if (res.ok) {
      alert('Disconnected successfully');
      window.location.reload();
    } else {
      const data = await res.json();
      console.log('Error data:', data);
      alert(`Failed to disconnect: ${data.message || 'Unknown error'}`);
    }
  } catch (err) {
    console.error('Disconnect error:', err);
    alert('Failed to disconnect');
  } finally {
    setDisconnecting(null);
  }
};

  return (
    <section className="section has-background-light">
      <Link to="/connections/requests" className="button connect is-link is-light">
        View connection requests
      </Link>
      <UserHeader />

      <div className="container">
        <div className="recommendations-container">
          <h1 className="title has-text-centered">Your connections</h1>
        </div>

        <div className="container-err">
        {loading && <p>Loading connections…</p>}
        {error && <p className="error-text">No connections found.</p>}
        {!loading && !error && data.length === 0 && (
          <p>No connections found.</p>
        )}
        </div>
        

        {!loading && !error && data.length > 0 && (
          <div className="columns is-multiline">
            {data.map((p) => {
              return (
                <div className="column is-12" key={p.id}>
                  <article className="box">
                    <h2 className="title is-5">{p.name}</h2>
                    <p className="is-size-6 has-text-grey">{p.addressCity}</p>
                    <p className="mt-2">
                      <strong>Child:</strong> {p.child.name}
                      <br />
                      <strong>Age:</strong> {p.child.ageYears}
                    </p>
                    <div className="buttons-container mt-3">
                      <button 
                        className="button con is-danger"
                        onClick={() => handleDisconnect(p.id)}
                        disabled={disconnecting === p.id}
                      >
                        {disconnecting === p.id ? 'Disconnecting...' : 'Disconnect'}
                      </button>
                      <button 
                        className="button con is-success"
                        onClick={() => openChat(p.id)}
                        disabled={chatLoading === p.id}
                      >
                        {chatLoading === p.id ? 'Opening...' : 'Chat'}
                      </button>
                    </div>
                  </article>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </section>
  );
}