import brush from "/brush_icon.svg";
import db from "/db_icon.svg";
import sun from "/sun_icon.svg";
import moon from "/moon_icon.svg";
import { useState } from "react";

export default function Settings() {
  const [activeTab, setActiveTab] = useState("db");
  return (
    <div className="settings">
      <div className="settings__overlay">
        <div className="settings__window">
          <div className="settings__buttons">
            <div
              className="settings__button"
              onClick={() => setActiveTab("db")}
            >
              <img src={db} alt="Настройки базы данных" />
              <p>База данных</p>
            </div>

            <div
              className="settings__button"
              onClick={() => setActiveTab("theme")}
            >
              <img src={brush} alt="Настройки внешнего вида" />
              <p>Внешний вид</p>
            </div>
          </div>
          <div className="settings__divider"></div>
          <div className="settings__parameters">
            {activeTab === "db" && (
              <>
                <input type="text" placeholder="логин" />
                <input type="text" placeholder="пароль" />
                <input type="text" placeholder="путь" />
              </>
            )}

            {activeTab === "theme" && (
              <div className="settings__themes">
                <div>
                  <img src={sun} alt="Светлая тема" />
                </div>
                <div>
                  <img src={moon} alt="Темная тема" />
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
