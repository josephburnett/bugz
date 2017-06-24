(ns client.core
  (:require-macros [cljs.core.async.macros :refer [go-loop]]
                   [cljs.core.async.macros :refer [go]])
  (:require [om.core :as om :include-macros true]
            [om.dom :as dom :include-macros true]
            [cljs.core.async :refer [chan put! <! >!]]
            [chord.client :refer [ws-ch]]))

(enable-console-print!)

(defonce app-state (atom {:owner "joe"
                          :world-view {}}))

(def ant "⚈")
(def color-dirt "#d3b383")
(def color-colony "#3d2501")
(def color-ant "#b50c03")

(defn cell-view [cell _]
  (reify om/IRender
    (render [_]
      (dom/td #js {:style #js {:width "20px" :height "20px"
                               :background (if (get cell "Colony") color-colony color-dirt)
                               :color color-ant}}
              (cond
                (nil? (get cell "Object")) (dom/div nil "")
                :default (condp = (get-in cell ["Object" "Direction"])
                           [1,0] (dom/div #js {:className "right"} ant)
                           [1,-1] (dom/div #js {:className "down-right"} ant)
                           [0,-1] (dom/div #js {:className "down"} ant)
                           [-1,-1] (dom/div #js {:className "down-left"} ant)
                           [-1,0] (dom/div #js {:className "left"} ant)
                           [-1,1] (dom/div #js {:className "up-left"} ant)
                           [0,1] (dom/div #js {:className "up"} ant)
                           [1,1] (dom/div #js {:className "up-right"} ant)))))))

(defn row-view [row _]
  (reify om/IRender
    (render [_]
      (apply dom/tr nil (om/build-all cell-view row)))))

(defn world-view [world _]
  (reify om/IRender
    (render [_]
      (if-not (contains? world "Points")
        (dom/div nil (dom/h2 nil "Loading..."))
        (let [rows (get world "Points")]
          (dom/div nil (dom/table nil (apply dom/tbody nil (om/build-all row-view rows)))))))))

(om/root
 (fn [data owner]
   (reify
     om/IRender
     (render [_]
       (dom/div nil
                (om/build world-view (:world-view data))))
     om/IDidMount
     (did-mount [_]
       (go
         (let [addr (str "ws://" (.-hostname js/location) ":8080/ws")
               {:keys [ws-channel error]} (<! (ws-ch addr {:format :json}))]
           (if error
             (print error)
             (do
               (set! (.-onkeydown js/document.body)
                     (fn [e]
                       (when (= " " (.-key e))
                         (go (>! ws-channel (clj->js {"Type" "ui-produce"
                                                      "Event" {"Owner" "joe"}}))))))
               (go-loop []
                 (let [{:keys [message error]} (<! ws-channel)]
                   (if (nil? message)
                     (print "channel closed")
                     (do
                       (when (= "view-update" (get message "Type"))
                         (om/transact! data #(assoc-in % [:world-view] (get-in message ["Event" "WorldView"]))))
                       (when-not error (recur)))))))))))))
  app-state
  {:target (. js/document (getElementById "app"))})
