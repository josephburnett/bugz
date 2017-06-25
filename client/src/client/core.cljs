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
(def color-dirt "#d3b383")
(def color-colony "#3d2501")
(def color-phermone "#ce9c52")
(def color-my-ant "#b50c03")
(def color-their-ant "#000")

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

(defn cell-view [cell _]
  (reify om/IRender
    (render [_]
      (dom/td #js {:onClick #(do
                               (om/update! cell ["Phermone"] (not (get cell "Phermone")))
                               (go (>! server-channel {"Type" "ui-phermone"
                                                       "Event" {"Point" (get cell "Point")
                                                                "State" (not (get cell "Phermone"))}})))
                   :style #js {:width "20px" :height "20px"
                               :background (cond
                                             (get cell "Colony") color-colony
                                             (get cell "Phermone") color-phermone
                                             :default color-dirt)}}
              (cond
                (nil? (get cell "Object")) (dom/div nil "")
                :default (let [style {:style {:color (if (get-in cell ["Object" "Mine"]) color-my-ant color-their-ant)}}]
                           (condp = (get-in cell ["Object" "Direction"])
                             [1,0] (dom/div (clj->js (merge style {:className "right"})) ant)
                             [1,-1] (dom/div (clj->js (merge style {:className "down-right"})) ant)
                             [0,-1] (dom/div (clj->js (merge style {:className "down"})) ant)
                             [-1,-1] (dom/div (clj->js (merge style {:className "down-left"})) ant)
                             [-1,0] (dom/div (clj->js (merge style {:className "left"})) ant)
                             [-1,1] (dom/div (clj->js (merge style {:className "up-left"})) ant)
                             [0,1] (dom/div (clj->js (merge style {:className "up"})) ant)
                             [1,1] (dom/div (clj->js (merge style {:className "up-right"})) ant))))))))

(defn row-view [row _]
  (reify om/IRender
    (render [_]
      (apply dom/tr nil (om/build-all cell-view row)))))

(defn world-view [world _]
  (reify
    om/IRender
    (render [_]
      (if-not (contains? world "Points")
        (show-text "Loading...")
        (let [rows (get world "Points")]
          (dom/div nil (dom/table nil (apply dom/tbody nil (om/build-all row-view rows)))))))
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
