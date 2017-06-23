(ns client.core
  (:require-macros [cljs.core.async.macros :refer [go-loop]]
                   [cljs.core.async.macros :refer [go]])
  (:require [om.core :as om :include-macros true]
            [om.dom :as dom :include-macros true]
            [cljs.core.async :refer [chan put! <! >!]]
            [chord.client :refer [ws-ch]]))

(enable-console-print!)

(println "This text is printed from src/client/core.cljs. Go ahead and edit it and see reloading in action.")

;; define your app data so that it doesn't get over-written on reload

(defonce app-state (atom {:owner "joe"
                          :world-view {}}))

(defn cell-view [cell _]
  (reify om/IRender
    (render [_]
      (dom/td #js {:style #js {:width "20px" :height "20px"}}
             (if (nil? (get cell "Object")) "" "A")))))

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
                (dom/h1 nil (str "Owner: "(:owner data)))
                (om/build world-view (:world-view data))))
     om/IDidMount
     (did-mount [_]
       (go
         (let [addr (str "ws://" (.-hostname js/location) ":8080/ws")
               {:keys [ws-channel error]} (<! (ws-ch addr {:format :json}))]
           (if error
             (print error)
             (go-loop []
               (print "listening")
               (let [{:keys [message error]} (<! ws-channel)]
                 (if (nil? message)
                   (print "channel closed")
                   (do
                     (print "got message")
                     (when (= "view-update" (get message "Type"))
                       (print "message is a view-update")
                       (om/transact! data #(assoc-in % [:world-view] (get-in message ["Event" "WorldView"]))))
                     (when-not error (recur))))))))))))
  app-state
  {:target (. js/document (getElementById "app"))})

(defn on-js-reload []
  ;; optionally touch your app-state to force rerendering depending on
  ;; your application
  ;; (swap! app-state update-in [:__figwheel_counter] inc)
)
