import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;
import '../models/clothing_item.dart';
import '../models/outfit.dart';
import '../models/outfit_recommendation.dart';

class ApiService {
  static const String baseUrl = 'http://localhost:8080/api/v1';
  
  // Device ID for authentication
  String deviceId = 'mobile-device-${DateTime.now().millisecondsSinceEpoch}';
  String? userId;

  Map<String, String> get headers => {
    'Content-Type': 'application/json',
    'X-Device-ID': deviceId,
  };

  // ========== USER ==========
  
  Future<Map<String, dynamic>> getCurrentUser() async {
    final response = await http.get(
      Uri.parse('$baseUrl/users/me'),
      headers: headers,
    );
    
    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      userId = data['id'];
      return data;
    }
    throw Exception('Failed to get user: ${response.statusCode}');
  }

  // ========== CLOTHING ITEMS ==========
  
  Future<List<ClothingItem>> getClothingItems() async {
    final response = await http.get(
      Uri.parse('$baseUrl/clothing'),
      headers: headers,
    );
    
    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      final items = data['items'] as List? ?? [];
      return items.map((e) => ClothingItem.fromJson(e)).toList();
    }
    throw Exception('Failed to get clothing items: ${response.statusCode}');
  }

  Future<ClothingItem> createClothingItem(ClothingItem item) async {
    final response = await http.post(
      Uri.parse('$baseUrl/clothing'),
      headers: headers,
      body: jsonEncode(item.toJson()),
    );
    
    if (response.statusCode == 201) {
      return ClothingItem.fromJson(jsonDecode(response.body));
    }
    throw Exception('Failed to create clothing item: ${response.statusCode}');
  }

  // ========== OUTFITS ==========
  
  Future<List<Outfit>> getOutfits(String userId) async {
    final response = await http.get(
      Uri.parse('$baseUrl/users/$userId/outfits'),
      headers: headers,
    );
    
    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      final outfits = data['outfits'] as List? ?? [];
      return outfits.map((e) => Outfit.fromJson(e)).toList();
    }
    throw Exception('Failed to get outfits: ${response.statusCode}');
  }

  Future<Outfit> createOutfit(Outfit outfit) async {
    final response = await http.post(
      Uri.parse('$baseUrl/outfits'),
      headers: headers,
      body: jsonEncode({
        'name': outfit.name,
        'description': outfit.description,
        'style': outfit.style,
        'occasion': outfit.occasion,
        'season': outfit.season,
        'weather': outfit.weather,
      }),
    );
    
    if (response.statusCode == 201) {
      return Outfit.fromJson(jsonDecode(response.body));
    }
    throw Exception('Failed to create outfit: ${response.statusCode}');
  }

  // ========== AI RECOMMENDATIONS ==========
  
  Future<List<OutfitRecommendation>> getRecommendations(
    String userId, {
    RecommendationRequest? request,
  }) async {
    final response = await http.post(
      Uri.parse('$baseUrl/users/$userId/outfits/recommendations'),
      headers: headers,
      body: jsonEncode(request?.toJson() ?? {}),
    );
    
    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      final recommendations = data['recommendations'] as List? ?? [];
      return recommendations.map((e) => OutfitRecommendation.fromJson(e)).toList();
    }
    throw Exception('Failed to get recommendations: ${response.statusCode}');
  }

  // ========== UPLOAD ==========
  
  Future<String> uploadImage(File imageFile) async {
    final request = http.MultipartRequest(
      'POST',
      Uri.parse('$baseUrl/upload'),
    );
    
    request.headers['X-Device-ID'] = deviceId;
    request.files.add(await http.MultipartFile.fromPath(
      'image',
      imageFile.path,
    ));

    final response = await request.send();
    
    if (response.statusCode == 200) {
      final responseData = await response.stream.bytesToString();
      final data = jsonDecode(responseData);
      return data['url'] ?? '';
    }
    throw Exception('Failed to upload image: ${response.statusCode}');
  }

  // ========== COLOR THEORY ==========
  
  Future<Map<String, dynamic>> getColorSeasons() async {
    final response = await http.get(
      Uri.parse('$baseUrl/color-theory/seasons'),
      headers: headers,
    );
    
    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    }
    throw Exception('Failed to get color seasons: ${response.statusCode}');
  }

  Future<Map<String, dynamic>> analyzeColorHarmony(List<String> colors) async {
    final response = await http.post(
      Uri.parse('$baseUrl/color-theory/analyze-harmony'),
      headers: headers,
      body: jsonEncode({'colors': colors}),
    );
    
    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    }
    throw Exception('Failed to analyze harmony: ${response.statusCode}');
  }
}
