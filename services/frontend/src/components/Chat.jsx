import { useState, useEffect, useRef } from "react";

import micro from "/micro_icon.svg";
import send from "/send_icon.svg";
import settings from "/settings_icon.svg";
import Settings from "./Settings";

export default function Chat() {
  const [text, setText] = useState("");
  const [listening, setListening] = useState(false);
  const [loading, setLoading] = useState(false);
  const [showSettings, setShowSettings] = useState(false);

  const [messages, setMessages] = useState([]);

  const messagesRef = useRef(null);
  const wsRef = useRef(null);
  const pingRef = useRef(null);
  const recognitionRef = useRef(null);

  const formatTime = (ts) => {
    if (ts) {
      const date = new Date(ts);
      if (!isNaN(date)) {
        return date.toLocaleTimeString([], {
          hour: "2-digit",
          minute: "2-digit",
        });
      }
    }
    return new Date().toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  useEffect(() => {
    if (messagesRef.current) {
      messagesRef.current.scrollTop = messagesRef.current.scrollHeight;
    }
  }, [messages, loading]);

  const parseMaps = (raw) => {
    const blocks = raw.match(/map\[(.*?)\]/g) || [];

    return blocks.map((block) => {
      const content = block.replace(/map\[|\]/g, "");

      const entries = content.match(/(\w+):(<nil>|[^ ]+)/g) || [];

      const obj = {};

      entries.forEach((entry) => {
        const [key, value] = entry.split(":");
        obj[key] = value === "<nil>" ? null : value;
      });

      return obj;
    });
  };

  const addMessage = (sender, text, created_at) => {
    setMessages((prev) => [
      ...prev,
      { sender, text, created_at, id: Date.now() + Math.random() },
    ]);
  };

  useEffect(() => {
    const ws = new WebSocket("wss://higu.su/ws");
    wsRef.current = ws;

    ws.onopen = () => {
      addMessage("system", "✅ Подключено к серверу");
      startPing();
    };

    ws.onclose = () => {
      addMessage("system", "❌ Соединение закрыто");
      clearInterval(pingRef.current);
    };

    ws.onerror = (error) => {
      console.error(error);
      addMessage("system", "⚠️ Ошибка соединения");
    };

    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data);

      if (msg.type === "pong") return;

      if (msg.username === "AI") {
        setLoading(false);

        const parsed = parseMaps(msg.text);

        const isTable = parsed.length > 0;

        addMessage("AI", isTable ? parsed : msg.text, msg.created_at);

        return;
      }

      addMessage(msg.username, msg.text, msg.created_at);
    };

    return () => {
      clearInterval(pingRef.current);
      ws.close();
    };
  }, []);

  const startPing = () => {
    pingRef.current = setInterval(() => {
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({ type: "ping" }));
      }
    }, 30000);
  };

  const handleSend = () => {
    if (!text.trim()) return;

    const ws = wsRef.current;
    if (!ws || ws.readyState !== WebSocket.OPEN) return;

    ws.send(JSON.stringify({ text }));

    setText("");

    setLoading(true);
  };

  const startListening = () => {
    const SpeechRecognition =
      window.SpeechRecognition || window.webkitSpeechRecognition;

    if (!SpeechRecognition) {
      alert("Браузер не поддерживает распознавание речи");
      return;
    }

    if (!recognitionRef.current) {
      const recognition = new SpeechRecognition();

      recognition.lang = "ru-RU";
      recognition.continuous = false;
      recognition.interimResults = false;

      recognition.onresult = (event) => {
        setText(event.results[0][0].transcript);
      };

      recognition.onend = () => setListening(false);

      recognitionRef.current = recognition;
    }

    setListening(true);
    recognitionRef.current.start();
  };

  return (
    <div className="chat">
      {showSettings && <Settings onClose={() => setShowSettings(false)} />}
      <div className="chat__container">
        <div className="chat__header">
          <h2 className="chat__title">Чат</h2>
          <img
            src={settings}
            alt="Настройки"
            onClick={() => setShowSettings(true)}
          />
        </div>

        <div className="chat__body" ref={messagesRef}>
          <div className="chat__messages">
            {messages.map((msg) => (
              <div
                key={msg.id}
                className={`chat__message chat__message--${
                  msg.sender === "User" ? "sent" : "received"
                }`}
              >
                <div className="chat__bubble">
                  {Array.isArray(msg.text) ? (
                    <div className="chat__tableWrapper">
                      <table className="chat__table">
                        <thead>
                          <tr>
                            {Object.keys(msg.text[0] || {}).map((key) => (
                              <th key={key}>{key}</th>
                            ))}
                          </tr>
                        </thead>

                        <tbody>
                          {msg.text.map((row, i) => (
                            <tr key={i}>
                              {Object.values(row).map((val, j) => (
                                <td key={j}>{val ?? "—"}</td>
                              ))}
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  ) : (
                    <div className="chat__text">{msg.text}</div>
                  )}

                  <div className="chat__time">{formatTime(msg.created_at)}</div>
                </div>
              </div>
            ))}

            {loading && (
              <div className="chat__message chat__message--received">
                <div className="chat__typing">
                  <div className="chat__dot"></div>
                  <div className="chat__dot"></div>
                  <div className="chat__dot"></div>
                </div>
              </div>
            )}
          </div>
        </div>

        <div className="chat__footer">
          <input
            className="chat__input"
            type="text"
            placeholder="Напишите сообщение"
            value={text}
            onChange={(e) => setText(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSend()}
          />

          <button
            className={`chat__button ${listening ? "chat__button--active" : ""}`}
            onClick={startListening}
          >
            <img src={micro} alt="voice" />
          </button>

          {text && (
            <button className="chat__button" onClick={handleSend}>
              <img src={send} alt="send" />
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
