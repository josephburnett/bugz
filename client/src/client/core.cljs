(ns client.core
  (:require [om.core :as om :include-macros true]
            [om.dom :as dom :include-macros true]
            [chord.client :refer [ws-ch]]))

(enable-console-print!)

(println "This text is printed from src/client/core.cljs. Go ahead and edit it and see reloading in action.")

;; define your app data so that it doesn't get over-written on reload

(defonce app-state (atom {:owner "joe"
                          :world-view {}}))

(defn cell-view [cell _]
  (reify om/IRender
    (render [_]
      (dom/td nil
             (if (nil? (:Object cell)) "0" "1")))))

(defn row-view [row _]
  (reify om/IRender
    (render [_]
      (apply dom/tr nil (om/build-all cell-view row)))))

(defn world-view [world _]
  (reify om/IRender
    (render [_]
      (if-not (contains? world :PointsView)
        (dom/div nil (dom/h2 nil "Loading..."))
        (let [points (:PointsView world)]
          (dom/div nil (dom/table nil (apply dom/tbody nil (om/build-all row-view world)))))))))

(om/root
  (fn [data owner]
    (reify om/IRender
      (render [_]
        (dom/div nil
                 (dom/h1 nil (str "Owner: "(:owner data)))
                 (om/build world-view (:world-view data))))))
  app-state
  {:target (. js/document (getElementById "app"))})

(defn on-js-reload []
  ;; optionally touch your app-state to force rerendering depending on
  ;; your application
  ;; (swap! app-state update-in [:__figwheel_counter] inc)
)
