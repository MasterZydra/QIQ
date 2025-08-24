<?php
$messages = [
    "Put your content here.",
    "Start editing this page with your own information.",
    "Add your custom message here.",
    "Replace this placeholder with your content.",
    "Your content goes here!"
];

$message = $messages[array_rand($messages)];

echo "<!DOCTYPE html>
<html>
<head>
    <title>QIQ Placeholder</title>
    <style>
        body { font-family: Arial, sans-serif; background: #f0f0f0; text-align: center; margin-top: 10%; }
        .box { background: #fff; padding: 40px; border-radius: 10px; display: inline-block; box-shadow: 0 2px 8px rgba(0,0,0,0.1);}
        h1 { color: #333; }
    </style>
</head>
<body>
    <div>
        <h1>QIQ</h1>
        <img src='Rabbit.svg' style='max-width: 10em; margin: 2em;'/>
    </div>
    <div class='box'>
        <h1>$message</h1>
    </div>
</body>
</html>";
?>