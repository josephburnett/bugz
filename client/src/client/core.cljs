(ns client.core
  (:require-macros [cljs.core.async.macros :refer [go-loop]]
                   [cljs.core.async.macros :refer [go]])
  (:require [om.core :as om :include-macros true]
            [om.dom :as dom :include-macros true]
            [cljs.core.async :refer [chan put! <! >!]]
            [chord.client :refer [ws-ch]]))

(enable-console-print!)

(defonce app-state (atom {:owner ""
                          :world-view {}}))

(defonce server-channel (chan))
(defonce client-channel (chan))

(def ant "⚈")
(def fruit "♥")
(def phermone "•")
(def color-dirt {0 "#d3b383"
                 1 "#c6a471"
                 2 "#c19553"
                 3 "#579657"})
(def color-colony "#3d2501")
(def color-phermone "#505ffc")
(def color-my-ant "#b50c03")
(def color-their-ant "#000")
(def color-friend "#1f6d1f")
(def color-enemy "#9e1914")
(def color-fruit "#0b8c0b")

(defn show-text [text]
  (dom/div nil (dom/h2 nil text)))

(defn connect [owner]
  (go
    (let [addr (str "ws://" (.-hostname js/location) ":8080/ws/owner/" (js/encodeURI owner))
          {:keys [ws-channel error]} (<! (ws-ch addr {:format :json}))]
      (if error
        (print error)
        (do
          (go-loop []
            (let [msg (<! server-channel)]
              (if (nil? msg)
                (print "server channel closed")
                (do
                  (let [owned-msg (clj->js (assoc-in msg ["Event" "Owner"] owner))]
                    (>! ws-channel owned-msg)
                    (recur))))))
          (go-loop []
            (let [{:keys [message error]} (<! ws-channel)]
              (if (nil? message)
                (print "websocket closed")
                (do
                  (>! client-channel message)
                  (recur))))))))))

(defn show-ant [cell]
  (let [style {:style {:color (if (get-in cell ["Object" "Mine"]) color-my-ant color-their-ant)
                       :position "absolute"
                       :top "0px"
                       :right "5px"}}]
    (condp = (get-in cell ["Object" "Direction"])
      [1,0] (dom/div (clj->js (merge style {:className "right"})) ant)
      [1,-1] (dom/div (clj->js (merge style {:className "down-right"})) ant)
      [0,-1] (dom/div (clj->js (merge style {:className "down"})) ant)
      [-1,-1] (dom/div (clj->js (merge style {:className "down-left"})) ant)
      [-1,0] (dom/div (clj->js (merge style {:className "left"})) ant)
      [-1,1] (dom/div (clj->js (merge style {:className "up-left"})) ant)
      [0,1] (dom/div (clj->js (merge style {:className "up"})) ant)
      [1,1] (dom/div (clj->js (merge style {:className "up-right"})) ant))))

(defn show-fruit [cell]
  (dom/div #js {:style #js {:fontSize "25px"
                            :color color-fruit
                            :position "absolute"
                            :top "-5px"
                            :right "3px"}}
           fruit))

(defn cell-view [cell _]
  (reify om/IRender
    (render [_]
      (dom/td #js {:onClick #(do
                               (om/update! cell ["Phermone"] (not (get cell "Phermone")))
                               (go (>! server-channel {"Type" "ui-phermone"
                                                       "Event" {"Point" (get cell "Point")
                                                                "State" (not (get cell "Phermone"))}})))
                   :style #js {:width "20px" :height "20px"
                               :textAlign "center"
                               :position "relative"
                               :background (cond
                                             (get cell "Colony") color-colony
                                             :default (get color-dirt (get cell "Soil")))}}
              (when (get cell "Phermone") (dom/div #js {:style #js {:position "relative"
                                                                    :color color-phermone
                                                                    :zIndex "100"}} phermone))
              (when-not (nil? (get cell "Object"))
                (if (= "ant" (get-in cell ["Object" "Type"]))
                  (show-ant cell)
                  (show-fruit cell)))))))

(defn row-view [row _]
  (reify om/IRender
    (render [_]
      (apply dom/tr nil (om/build-all cell-view row)))))

(defn friend-view [friends _]
  (reify
    om/IRender
    (render [_]
      (if (= 0 (count friends))
        (dom/span nil "no other colonies")
        (apply dom/ul nil
               (map (fn [f]
                      (let [is-friend (get friends f)]
                        (dom/li #js {:onClick #(go (>! server-channel {"Type" "ui-friend"
                                                                       "Event" {"Friend" f
                                                                                "State" (not is-friend)}}))}
                                (dom/span #js {:style #js {:color (if is-friend color-friend color-enemy)}} f))))
                    (keys friends)))))))

(defn world-view [world _]
  (reify
    om/IRender
    (render [_]
      (if-not (contains? world "Points")
        (show-text "Loading...")
        (let [rows (get world "Points")]
          (dom/div nil
                   (dom/table nil (apply dom/tbody nil (om/build-all row-view rows)))
                   (om/build friend-view (get world "Friends"))))))
    om/IDidMount
    (did-mount [_]
      (set! (.-onkeydown js/document.body)
            (fn [e]
              (when (= " " (.-key e))
                (go (>! server-channel {"Type" "ui-produce"
                                        "Event" {}})))))
      (go-loop []
        (let [msg (<! client-channel)]
          (when (= "view-update" (get msg "Type"))
            (om/update! world (get-in msg ["Event" "WorldView"]))))
        (recur)))))

(defn owner-selection [data owner]
  (reify
    om/IInitState
    (init-state [_]
      {:owner ""})
    om/IRenderState
    (render-state [_ state]
      (dom/div nil
               (dom/h1 nil "Colony name:")
               (show-text (:owner state))))
    om/IDidMount
    (did-mount [_]
      (set! (.-onkeydown js/document.body)
            (fn [e] (let [k (.-key e)
                          o (om/get-state owner :owner)]
                      (cond
                        (and (= "Enter" k)
                             (not (= "" o)))
                        (do
                          (om/update! data [:owner] o)
                          (connect o))
                        (= "Backspace" k)
                        (when-not (= "" o)
                          (om/set-state! owner [:owner] (subs o 0 (dec (count o)))))
                        (< 100 (count o))
                        (print "owner too long")
                        (re-matches #"[a-z\-]" k)
                        (om/set-state! owner [:owner] (str o k)))))))))

(om/root
 (fn [data owner]
   (reify
     om/IRender
     (render [_]
       (dom/div nil
                (if (empty? (:owner data))
                  (om/build owner-selection data)
                  (om/build world-view (:world-view data)))))))
 app-state
 {:target (. js/document (getElementById "app"))})
