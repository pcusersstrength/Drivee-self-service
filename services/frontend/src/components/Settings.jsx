import brush from "/brush_icon.svg";
import db from "/db_icon.svg";
import sun from "/sun_icon.svg";
import moon from "/moon_icon.svg";
import { useState } from "react";

export default function Settings({ onClose }) {
  const [activeTab, setActiveTab] = useState("db");

  return (
    <div className="settings">
      <div className="settings__overlay" onClick={onClose}>
        <div className="settings__window" onClick={(e) => e.stopPropagation()}>
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
                <div className="settings__parameter">
                  <label>Логин</label>
                  <input type="text" />
                </div>

                <div className="settings__parameter">
                  <label>Пароль</label>
                  <input type="text" />
                </div>

                <div className="settings__parameter">
                  <label>Путь</label>
                  <input type="text" />
                </div>
              </>
            )}

            {activeTab === "theme" && (
              <div className="settings__themes">
                <div className="settings__theme-button">
                  <img src={sun} alt="Светлая тема" />
                </div>

                <div className="settings__theme-button settings__theme-button--active">
                  <img src={moon} alt="Темная тема" />
                </div>
              </div>
            )}
            <button className="settings__save">Сохранить изменения</button>
          </div>
        </div>
      </div>
    </div>
  );
}
